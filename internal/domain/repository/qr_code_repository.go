package repository

import "github.com/stempo/backend/internal/domain/entity"

type QRCodeRepository interface {
	Create(qrCode *entity.QRCode) error
	FindByCode(code string) (*entity.QRCode, error)
	FindByID(id uint) (*entity.QRCode, error)
	FindByBonusProgramID(programID uint) ([]entity.QRCode, error)
	FindByBusinessID(businessID uint) ([]entity.QRCode, error)
	Update(qrCode *entity.QRCode) error
	Delete(id uint) error
	DeleteByBonusProgramID(programID uint) error
}
