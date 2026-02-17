package usecase

import (
	"errors"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
)

type BusinessUsecase interface {
	Create(business *entity.Business) error
	GetByID(id uint) (*entity.Business, error)
	GetByOwnerID(ownerID uint) (*entity.Business, error)
	GetFeatured() ([]entity.Business, error)
	GetByCategory(category string) ([]entity.Business, error)
	GetNewest(limit int) ([]entity.Business, error)
	GetAll(page, pageSize int) ([]entity.Business, int64, error)
	GetByStatus(status entity.BusinessStatus, page, pageSize int) ([]entity.Business, int64, error)
	Approve(id uint) error
	Reject(id uint, reason string) error
	Update(business *entity.Business, userID uint) error
	UpdateFields(id uint, updates map[string]interface{}) error
	Delete(id, userID uint) error
}

type businessUsecase struct {
	businessRepo repository.BusinessRepository
}

func NewBusinessUsecase(businessRepo repository.BusinessRepository) BusinessUsecase {
	return &businessUsecase{businessRepo: businessRepo}
}

func (u *businessUsecase) Create(business *entity.Business) error {
	return u.businessRepo.Create(business)
}

func (u *businessUsecase) GetByID(id uint) (*entity.Business, error) {
	return u.businessRepo.FindByID(id)
}

func (u *businessUsecase) GetByOwnerID(ownerID uint) (*entity.Business, error) {
	business, err := u.businessRepo.FindByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}
	// If business is nil (not found), return nil without error
	return business, nil
}

func (u *businessUsecase) GetFeatured() ([]entity.Business, error) {
	return u.businessRepo.FindFeatured()
}

func (u *businessUsecase) GetByCategory(category string) ([]entity.Business, error) {
	return u.businessRepo.FindByCategory(category)
}

func (u *businessUsecase) GetNewest(limit int) ([]entity.Business, error) {
	if limit < 1 {
		limit = 10
	}
	return u.businessRepo.FindNewest(limit)
}

func (u *businessUsecase) GetAll(page, pageSize int) ([]entity.Business, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return u.businessRepo.FindAll(page, pageSize)
}

func (u *businessUsecase) GetByStatus(status entity.BusinessStatus, page, pageSize int) ([]entity.Business, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return u.businessRepo.FindByStatus(status, page, pageSize)
}

func (u *businessUsecase) Approve(id uint) error {
	business, err := u.businessRepo.FindByID(id)
	if err != nil {
		return errors.New("business not found")
	}
	business.Status = entity.BusinessStatusApproved
	business.RejectionReason = nil
	return u.businessRepo.Update(business)
}

func (u *businessUsecase) Reject(id uint, reason string) error {
	business, err := u.businessRepo.FindByID(id)
	if err != nil {
		return errors.New("business not found")
	}
	business.Status = entity.BusinessStatusRejected
	reasonStr := reason
	business.RejectionReason = &reasonStr
	return u.businessRepo.Update(business)
}

func (u *businessUsecase) Update(business *entity.Business, userID uint) error {
	existing, err := u.businessRepo.FindByID(business.ID)
	if err != nil {
		return errors.New("business not found")
	}

	// If userID is 0, allow update (admin override)
	if userID != 0 && existing.OwnerID != userID {
		return errors.New("unauthorized")
	}

	// Preserve system-managed fields only if not admin (userID == 0 means admin)
	if userID != 0 {
		business.OwnerID = existing.OwnerID
		business.Status = existing.Status
		business.Rating = existing.Rating
	} else {
		// Admin can change status and rating, but preserve OwnerID and CreatedAt
		business.OwnerID = existing.OwnerID
		business.CreatedAt = existing.CreatedAt
	}
	// UpdatedAt will be set automatically by GORM
	
	return u.businessRepo.Update(business)
}

func (u *businessUsecase) UpdateFields(id uint, updates map[string]interface{}) error {
	return u.businessRepo.UpdateFields(id, updates)
}

func (u *businessUsecase) Delete(id, userID uint) error {
	business, err := u.businessRepo.FindByID(id)
	if err != nil {
		return errors.New("business not found")
	}

	if business.OwnerID != userID {
		return errors.New("unauthorized")
	}

	return u.businessRepo.Delete(id)
}
