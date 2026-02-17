package entity

import "time"

type QRCode struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Code          string     `gorm:"not null;uniqueIndex" json:"code"`
	Type          QRCodeType `gorm:"not null;type:varchar(20)" json:"type"`
	BonusProgramID uint      `gorm:"not null" json:"bonusProgramId"`
	BonusProgram  BonusProgram `gorm:"foreignKey:BonusProgramID" json:"bonusProgram,omitempty"`
	BusinessID    uint      `gorm:"not null" json:"businessId"`
	Business      Business  `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
	IsActive      bool       `gorm:"default:true" json:"isActive"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func (QRCode) TableName() string {
	return "qr_codes"
}
