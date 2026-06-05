package postgres

import (
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type bonusProgramRepository struct {
	db *gorm.DB
}

func NewBonusProgramRepository(db *gorm.DB) repository.BonusProgramRepository {
	return &bonusProgramRepository{db: db}
}

func (r *bonusProgramRepository) Create(program *entity.BonusProgram) error {
	return r.db.Create(program).Error
}

func (r *bonusProgramRepository) FindByID(id uint) (*entity.BonusProgram, error) {
	var program entity.BonusProgram
	err := r.db.Preload("Business").Preload("Business.Owner").First(&program, id).Error
	if err != nil {
		return nil, err
	}
	return &program, nil
}

// FindByBusinessID returns ALL programs for a business, regardless of status
// This is used by business owners to see all their programs (including pending)
func (r *bonusProgramRepository) FindByBusinessID(businessID uint) ([]entity.BonusProgram, error) {
	var programs []entity.BonusProgram
	// NO status filtering - return all programs
	err := r.db.Where("business_id = ?", businessID).
		Preload("Business").
		Order("created_at DESC").
		Find(&programs).Error
	return programs, err
}

// FindActiveByBusinessID returns only approved/active programs for public listing.
func (r *bonusProgramRepository) FindActiveByBusinessID(businessID uint) ([]entity.BonusProgram, error) {
	var programs []entity.BonusProgram
	err := r.db.Where("business_id = ? AND status IN ?", businessID,
		[]entity.BonusProgramStatus{entity.BonusProgramStatusApproved, entity.BonusProgramStatusActive}).
		Preload("Business").
		Order("created_at DESC").
		Find(&programs).Error
	return programs, err
}

func (r *bonusProgramRepository) FindByStatus(status entity.BonusProgramStatus, page, pageSize int) ([]entity.BonusProgram, int64, error) {
	var programs []entity.BonusProgram
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&entity.BonusProgram{}).Where("status = ?", status).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("status = ?", status).
		Preload("Business").
		Preload("Business.Owner").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&programs).Error
	return programs, total, err
}

// FindAll returns only approved/active programs for public listing
func (r *bonusProgramRepository) FindAll(page, pageSize int) ([]entity.BonusProgram, int64, error) {
	var programs []entity.BonusProgram
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&entity.BonusProgram{}).
		Where("status IN ?", []entity.BonusProgramStatus{entity.BonusProgramStatusApproved, entity.BonusProgramStatusActive}).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("status IN ?", []entity.BonusProgramStatus{entity.BonusProgramStatusApproved, entity.BonusProgramStatusActive}).
		Preload("Business").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&programs).Error
	return programs, total, err
}

func (r *bonusProgramRepository) Update(program *entity.BonusProgram) error {
	return r.db.Save(program).Error
}

func (r *bonusProgramRepository) Delete(id uint) error {
	return r.db.Delete(&entity.BonusProgram{}, id).Error
}
