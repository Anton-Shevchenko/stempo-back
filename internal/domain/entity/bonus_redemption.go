package entity

import "time"

type BonusRedemption struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CardID     uint       `gorm:"not null;index" json:"cardId"`
	Card       LoyaltyCard `gorm:"foreignKey:CardID" json:"card,omitempty"`
	Code       string     `gorm:"uniqueIndex;not null;type:varchar(255)" json:"code"`
	ExpiresAt  time.Time  `gorm:"not null" json:"expiresAt"`
	Used       bool       `gorm:"default:false" json:"used"`
	UsedAt     *time.Time `json:"usedAt,omitempty"`
	ScannedBy  *uint      `json:"scannedBy,omitempty"`
	ScannedByUser *User   `gorm:"foreignKey:ScannedBy" json:"scannedByUser,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

func (BonusRedemption) TableName() string {
	return "bonus_redemptions"
}
