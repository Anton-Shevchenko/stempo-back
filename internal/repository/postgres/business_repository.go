package postgres

import (
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type businessRepository struct {
	db *gorm.DB
}

func NewBusinessRepository(db *gorm.DB) repository.BusinessRepository {
	return &businessRepository{db: db}
}

func (r *businessRepository) Create(business *entity.Business) error {
	return r.db.Create(business).Error
}

func (r *businessRepository) FindByID(id uint) (*entity.Business, error) {
	var business entity.Business
	err := r.db.Preload("Owner").First(&business, id).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

func (r *businessRepository) FindByOwnerID(ownerID uint) (*entity.Business, error) {
	var business entity.Business
	err := r.db.Where("owner_id = ?", ownerID).First(&business).Error
	if err != nil {
		// Return nil, nil if record not found (not an error)
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &business, nil
}

func (r *businessRepository) FindFeatured() ([]entity.Business, error) {
	var businesses []entity.Business
	// Use RANDOM() for PostgreSQL to shuffle results
	err := r.db.Where("featured = ? AND status = ?", true, entity.BusinessStatusApproved).
		Preload("Owner").
		Order("RANDOM()").
		Find(&businesses).Error
	return businesses, err
}

func (r *businessRepository) FindByCategory(category string) ([]entity.Business, error) {
	var businesses []entity.Business
	// Return approved businesses for the category in random order
	err := r.db.Where("category = ? AND status = ?", category, entity.BusinessStatusApproved).
		Preload("Owner").
		Order("RANDOM()").
		Find(&businesses).Error
	return businesses, err
}

func (r *businessRepository) FindNewest(limit int) ([]entity.Business, error) {
	var businesses []entity.Business
	err := r.db.Where("status = ?", entity.BusinessStatusApproved).Order("created_at DESC").Limit(limit).Preload("Owner").Find(&businesses).Error
	return businesses, err
}

func (r *businessRepository) FindByStatus(status entity.BusinessStatus, page, pageSize int) ([]entity.Business, int64, error) {
	var businesses []entity.Business
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&entity.Business{}).Where("status = ?", status).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("status = ?", status).Preload("Owner").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&businesses).Error
	return businesses, total, err
}

func (r *businessRepository) FindAll(page, pageSize int) ([]entity.Business, int64, error) {
	var businesses []entity.Business
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&entity.Business{}).Where("status = ?", entity.BusinessStatusApproved).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("status = ?", entity.BusinessStatusApproved).Preload("Owner").Offset(offset).Limit(pageSize).Find(&businesses).Error
	return businesses, total, err
}

func (r *businessRepository) Update(business *entity.Business) error {
	return r.db.Save(business).Error
}

func (r *businessRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	// GORM's Updates method automatically converts camelCase to snake_case
	// and only updates non-zero values
	return r.db.Model(&entity.Business{}).Where("id = ?", id).Updates(updates).Error
}

func (r *businessRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Business{}, id).Error
}
