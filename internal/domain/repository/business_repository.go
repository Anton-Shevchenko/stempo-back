package repository

import "github.com/stempo/backend/internal/domain/entity"

type BusinessRepository interface {
	Create(business *entity.Business) error
	FindByID(id uint) (*entity.Business, error)
	FindByOwnerID(ownerID uint) (*entity.Business, error)
	FindFeatured() ([]entity.Business, error)
	FindByCategory(category string) ([]entity.Business, error)
	FindNewest(limit int) ([]entity.Business, error)
	FindAll(page, pageSize int) ([]entity.Business, int64, error)
	FindByStatus(status entity.BusinessStatus, page, pageSize int) ([]entity.Business, int64, error)
	Update(business *entity.Business) error
	UpdateFields(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}
