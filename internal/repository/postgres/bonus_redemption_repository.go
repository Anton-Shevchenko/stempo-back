package postgres

import (
	"time"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type bonusRedemptionRepository struct {
	db *gorm.DB
}

func NewBonusRedemptionRepository(db *gorm.DB) repository.BonusRedemptionRepository {
	return &bonusRedemptionRepository{db: db}
}

func (r *bonusRedemptionRepository) Create(redemption *entity.BonusRedemption) error {
	return r.db.Create(redemption).Error
}

func (r *bonusRedemptionRepository) FindByCode(code string) (*entity.BonusRedemption, error) {
	var redemption entity.BonusRedemption
	err := r.db.Where("code = ?", code).
		Preload("Card").
		Preload("Card.User").
		Preload("Card.Business").
		Preload("ScannedByUser").
		First(&redemption).Error
	if err != nil {
		return nil, err
	}
	return &redemption, nil
}

func (r *bonusRedemptionRepository) FindByCardID(cardID uint) ([]entity.BonusRedemption, error) {
	var redemptions []entity.BonusRedemption
	err := r.db.Where("card_id = ?", cardID).
		Preload("Card").
		Preload("ScannedByUser").
		Order("created_at DESC").
		Find(&redemptions).Error
	return redemptions, err
}

func (r *bonusRedemptionRepository) Update(redemption *entity.BonusRedemption) error {
	return r.db.Save(redemption).Error
}

func (r *bonusRedemptionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ? AND used = ?", time.Now(), false).
		Delete(&entity.BonusRedemption{}).Error
}
