package repository

import "github.com/stempo/backend/internal/domain/entity"

type LoyaltyCardRepository interface {
	Create(card *entity.LoyaltyCard) error
	FindByID(id uint) (*entity.LoyaltyCard, error)
	FindByUserID(userID uint) ([]entity.LoyaltyCard, error)
	FindByBusinessID(businessID uint) ([]entity.LoyaltyCard, error)
	FindByUserAndBusiness(userID, businessID uint) (*entity.LoyaltyCard, error)
	Update(card *entity.LoyaltyCard) error
	Delete(id uint) error
}
