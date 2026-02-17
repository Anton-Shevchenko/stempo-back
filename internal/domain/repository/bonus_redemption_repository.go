package repository

import "github.com/stempo/backend/internal/domain/entity"

type BonusRedemptionRepository interface {
	Create(redemption *entity.BonusRedemption) error
	FindByCode(code string) (*entity.BonusRedemption, error)
	FindByCardID(cardID uint) ([]entity.BonusRedemption, error)
	Update(redemption *entity.BonusRedemption) error
	DeleteExpired() error
}
