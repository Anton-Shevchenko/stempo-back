package usecase

import (
	"errors"
	"time"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"github.com/stempo/backend/pkg/utils"
)

type QRCodeUsecase interface {
	Generate(programID uint, qrType entity.QRCodeType, userID uint) (*entity.QRCode, error)
	Validate(code string) (*entity.QRCode, error)
	GetByProgramID(programID uint, userID uint) ([]entity.QRCode, error)
	GetByBusinessID(businessID uint, userID uint) ([]entity.QRCode, error)
	Delete(id uint, userID uint) error
}

type qrCodeUsecase struct {
	qrCodeRepo   repository.QRCodeRepository
	programRepo  repository.BonusProgramRepository
	businessRepo repository.BusinessRepository
}

func NewQRCodeUsecase(
	qrCodeRepo repository.QRCodeRepository,
	programRepo repository.BonusProgramRepository,
	businessRepo repository.BusinessRepository,
) QRCodeUsecase {
	return &qrCodeUsecase{
		qrCodeRepo:   qrCodeRepo,
		programRepo:  programRepo,
		businessRepo: businessRepo,
	}
}

func (u *qrCodeUsecase) Generate(programID uint, qrType entity.QRCodeType, userID uint) (*entity.QRCode, error) {
	// Verify program exists and user owns the business
	program, err := u.programRepo.FindByID(programID)
	if err != nil {
		return nil, errors.New("bonus program not found")
	}

	// Cannot generate QR for unapproved programs
	if program.Status != entity.BonusProgramStatusApproved && program.Status != entity.BonusProgramStatusActive {
		return nil, errors.New("cannot generate QR for unapproved bonus program")
	}

	business, err := u.businessRepo.FindByID(program.BusinessID)
	if err != nil || business == nil {
		return nil, errors.New("business not found")
	}

	if business.OwnerID != userID {
		return nil, errors.New("unauthorized")
	}

	// Generate unique code
	code, err := utils.GenerateUniqueCode()
	if err != nil {
		return nil, errors.New("failed to generate QR code")
	}

	// Set expiration for temporary QR codes (use program's QRExpirationMinutes: 1-10, default 3)
	var expiresAt *time.Time
	if qrType == entity.QRCodeTypeTemporary {
		minutes := program.QRExpirationMinutes
		if minutes < 1 {
			minutes = 3
		}
		if minutes > 10 {
			minutes = 10
		}
		exp := time.Now().Add(time.Duration(minutes) * time.Minute)
		expiresAt = &exp
	}

	qrCode := &entity.QRCode{
		Code:           code,
		Type:           qrType,
		BonusProgramID: programID,
		BusinessID:     program.BusinessID,
		ExpiresAt:      expiresAt,
		IsActive:       true,
	}

	if err := u.qrCodeRepo.Create(qrCode); err != nil {
		return nil, errors.New("failed to create QR code")
	}

	return qrCode, nil
}

func (u *qrCodeUsecase) Validate(code string) (*entity.QRCode, error) {
	qrCode, err := u.qrCodeRepo.FindByCode(code)
	if err != nil {
		return nil, errors.New("QR code not found")
	}

	if !qrCode.IsActive {
		return nil, errors.New("QR code is inactive")
	}

	// Check expiration for temporary QR codes
	if qrCode.Type == entity.QRCodeTypeTemporary && qrCode.ExpiresAt != nil {
		if time.Now().After(*qrCode.ExpiresAt) {
			return nil, errors.New("QR code has expired")
		}
	}

	return qrCode, nil
}

func (u *qrCodeUsecase) GetByProgramID(programID uint, userID uint) ([]entity.QRCode, error) {
	// Verify program exists and user owns the business
	program, err := u.programRepo.FindByID(programID)
	if err != nil {
		return nil, errors.New("bonus program not found")
	}

	business, err := u.businessRepo.FindByID(program.BusinessID)
	if err != nil || business == nil {
		return nil, errors.New("business not found")
	}

	if business.OwnerID != userID {
		return nil, errors.New("unauthorized")
	}

	return u.qrCodeRepo.FindByBonusProgramID(programID)
}

func (u *qrCodeUsecase) GetByBusinessID(businessID uint, userID uint) ([]entity.QRCode, error) {
	business, err := u.businessRepo.FindByID(businessID)
	if err != nil || business == nil {
		return nil, errors.New("business not found")
	}

	if business.OwnerID != userID {
		return nil, errors.New("unauthorized")
	}

	return u.qrCodeRepo.FindByBusinessID(businessID)
}

func (u *qrCodeUsecase) Delete(id uint, userID uint) error {
	qrCode, err := u.qrCodeRepo.FindByID(id)
	if err != nil {
		return errors.New("QR code not found")
	}

	business, err := u.businessRepo.FindByID(qrCode.BusinessID)
	if err != nil || business == nil {
		return errors.New("business not found")
	}

	if business.OwnerID != userID {
		return errors.New("unauthorized")
	}

	return u.qrCodeRepo.Delete(id)
}
