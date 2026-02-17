package postgres

import (
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type qrCodeRepository struct {
	db *gorm.DB
}

func NewQRCodeRepository(db *gorm.DB) repository.QRCodeRepository {
	return &qrCodeRepository{db: db}
}

func (r *qrCodeRepository) Create(qrCode *entity.QRCode) error {
	return r.db.Create(qrCode).Error
}

func (r *qrCodeRepository) FindByCode(code string) (*entity.QRCode, error) {
	var qrCode entity.QRCode
	err := r.db.Where("code = ?", code).
		Preload("BonusProgram").
		Preload("Business").
		First(&qrCode).Error
	if err != nil {
		return nil, err
	}
	return &qrCode, nil
}

func (r *qrCodeRepository) FindByID(id uint) (*entity.QRCode, error) {
	var qrCode entity.QRCode
	err := r.db.Preload("BonusProgram").
		Preload("Business").
		First(&qrCode, id).Error
	if err != nil {
		return nil, err
	}
	return &qrCode, nil
}

func (r *qrCodeRepository) FindByBonusProgramID(programID uint) ([]entity.QRCode, error) {
	var qrCodes []entity.QRCode
	err := r.db.Where("bonus_program_id = ?", programID).
		Preload("BonusProgram").
		Preload("Business").
		Order("created_at DESC").
		Find(&qrCodes).Error
	return qrCodes, err
}

func (r *qrCodeRepository) FindByBusinessID(businessID uint) ([]entity.QRCode, error) {
	var qrCodes []entity.QRCode
	err := r.db.Where("business_id = ?", businessID).
		Preload("BonusProgram").
		Preload("Business").
		Order("created_at DESC").
		Find(&qrCodes).Error
	return qrCodes, err
}

func (r *qrCodeRepository) Update(qrCode *entity.QRCode) error {
	return r.db.Save(qrCode).Error
}

func (r *qrCodeRepository) Delete(id uint) error {
	return r.db.Delete(&entity.QRCode{}, id).Error
}

func (r *qrCodeRepository) DeleteByBonusProgramID(programID uint) error {
	return r.db.Where("bonus_program_id = ?", programID).Delete(&entity.QRCode{}).Error
}
