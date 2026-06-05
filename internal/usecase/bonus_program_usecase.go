package usecase

import (
	"errors"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
)

type BonusProgramUsecase interface {
	Create(program *entity.BonusProgram, userID uint) error
	GetByID(id uint) (*entity.BonusProgram, error)
	GetByBusinessID(businessID uint) ([]entity.BonusProgram, error)
	GetActiveByBusinessID(businessID uint) ([]entity.BonusProgram, error)
	GetAll(page, pageSize int) ([]entity.BonusProgram, int64, error)
	GetByStatus(status entity.BonusProgramStatus, page, pageSize int) ([]entity.BonusProgram, int64, error)
	Approve(id uint) error
	Reject(id uint, reason string) error
	Update(program *entity.BonusProgram, userID uint) error
	Delete(id, userID uint) error
}

type bonusProgramUsecase struct {
	programRepo  repository.BonusProgramRepository
	businessRepo repository.BusinessRepository
}

func NewBonusProgramUsecase(
	programRepo repository.BonusProgramRepository,
	businessRepo repository.BusinessRepository,
) BonusProgramUsecase {
	return &bonusProgramUsecase{
		programRepo:  programRepo,
		businessRepo: businessRepo,
	}
}

// Helper function to verify business ownership
func (u *bonusProgramUsecase) verifyBusinessOwnership(businessID, userID uint) (*entity.Business, error) {
	business, err := u.businessRepo.FindByID(businessID)
	if err != nil {
		return nil, errors.New("business not found")
	}

	if business == nil {
		return nil, errors.New("business not found")
	}

	if business.OwnerID != userID {
		return nil, errors.New("unauthorized")
	}

	return business, nil
}

func (u *bonusProgramUsecase) Create(program *entity.BonusProgram, userID uint) error {
	if _, err := u.verifyBusinessOwnership(program.BusinessID, userID); err != nil {
		return err
	}

	// Always set status to pending - only admin can approve via admin panel
	program.Status = entity.BonusProgramStatusPending
	program.RejectionReason = nil

	// QR expiration for temporary QR: 1-10 minutes, default 3
	if program.QRExpirationMinutes < 1 || program.QRExpirationMinutes > 10 {
		program.QRExpirationMinutes = 3
	}

	return u.programRepo.Create(program)
}

func (u *bonusProgramUsecase) GetByID(id uint) (*entity.BonusProgram, error) {
	return u.programRepo.FindByID(id)
}

func (u *bonusProgramUsecase) GetByBusinessID(businessID uint) ([]entity.BonusProgram, error) {
	return u.programRepo.FindByBusinessID(businessID)
}

func (u *bonusProgramUsecase) GetActiveByBusinessID(businessID uint) ([]entity.BonusProgram, error) {
	return u.programRepo.FindActiveByBusinessID(businessID)
}

func (u *bonusProgramUsecase) GetAll(page, pageSize int) ([]entity.BonusProgram, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return u.programRepo.FindAll(page, pageSize)
}

func (u *bonusProgramUsecase) GetByStatus(status entity.BonusProgramStatus, page, pageSize int) ([]entity.BonusProgram, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return u.programRepo.FindByStatus(status, page, pageSize)
}

func (u *bonusProgramUsecase) Approve(id uint) error {
	program, err := u.programRepo.FindByID(id)
	if err != nil {
		return errors.New("bonus program not found")
	}

	program.Status = entity.BonusProgramStatusApproved
	program.RejectionReason = nil

	return u.programRepo.Update(program)
}

func (u *bonusProgramUsecase) Reject(id uint, reason string) error {
	program, err := u.programRepo.FindByID(id)
	if err != nil {
		return errors.New("bonus program not found")
	}

	program.Status = entity.BonusProgramStatusRejected
	reasonStr := reason
	program.RejectionReason = &reasonStr

	return u.programRepo.Update(program)
}

func (u *bonusProgramUsecase) Update(program *entity.BonusProgram, userID uint) error {
	existing, err := u.programRepo.FindByID(program.ID)
	if err != nil {
		return errors.New("program not found")
	}

	if _, err := u.verifyBusinessOwnership(existing.BusinessID, userID); err != nil {
		return err
	}

	// Prevent status changes - only admin can change status via Approve/Reject
	// Exception: if rejected, allow resubmission by setting back to pending
	if program.Status != existing.Status {
		if existing.Status == entity.BonusProgramStatusRejected && program.Status == entity.BonusProgramStatusPending {
			// Allow resubmission after rejection
			program.Status = entity.BonusProgramStatusPending
			program.RejectionReason = nil
		} else {
			// Preserve existing status for all other cases
			program.Status = existing.Status
		}
	}

	// QR expiration for temporary QR: 1-10 minutes, default 3
	if program.QRExpirationMinutes < 1 || program.QRExpirationMinutes > 10 {
		program.QRExpirationMinutes = 3
	}

	return u.programRepo.Update(program)
}

func (u *bonusProgramUsecase) Delete(id, userID uint) error {
	program, err := u.programRepo.FindByID(id)
	if err != nil {
		return errors.New("program not found")
	}

	if _, err := u.verifyBusinessOwnership(program.BusinessID, userID); err != nil {
		return err
	}

	return u.programRepo.Delete(id)
}
