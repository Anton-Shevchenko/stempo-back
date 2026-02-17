package postgres

import (
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type loyaltyCardRepository struct {
	db *gorm.DB
}

func NewLoyaltyCardRepository(db *gorm.DB) repository.LoyaltyCardRepository {
	return &loyaltyCardRepository{db: db}
}

func (r *loyaltyCardRepository) Create(card *entity.LoyaltyCard) error {
	return r.db.Create(card).Error
}

func (r *loyaltyCardRepository) FindByID(id uint) (*entity.LoyaltyCard, error) {
	var card entity.LoyaltyCard
	err := r.db.Preload("User").Preload("Business").Preload("Business.Owner").First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *loyaltyCardRepository) FindByUserID(userID uint) ([]entity.LoyaltyCard, error) {
	var cards []entity.LoyaltyCard
	err := r.db.Where("user_id = ?", userID).Preload("Business").Find(&cards).Error
	return cards, err
}

func (r *loyaltyCardRepository) FindByBusinessID(businessID uint) ([]entity.LoyaltyCard, error) {
	var cards []entity.LoyaltyCard
	err := r.db.Where("business_id = ?", businessID).Preload("User").Find(&cards).Error
	return cards, err
}

func (r *loyaltyCardRepository) FindByUserAndBusiness(userID, businessID uint) (*entity.LoyaltyCard, error) {
	var card entity.LoyaltyCard
	err := r.db.Where("user_id = ? AND business_id = ?", userID, businessID).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *loyaltyCardRepository) Update(card *entity.LoyaltyCard) error {
	return r.db.Save(card).Error
}

func (r *loyaltyCardRepository) Delete(id uint) error {
	return r.db.Delete(&entity.LoyaltyCard{}, id).Error
}
