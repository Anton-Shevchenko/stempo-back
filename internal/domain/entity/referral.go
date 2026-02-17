package entity

import "time"

type ReferralStatus string

const (
	ReferralStatusPending  ReferralStatus = "pending"
	ReferralStatusCompleted ReferralStatus = "completed"
	ReferralStatusRewarded ReferralStatus = "rewarded"
)

type Referral struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ReferrerID    uint           `gorm:"not null" json:"referrerId"`
	Referrer      User           `gorm:"foreignKey:ReferrerID" json:"referrer,omitempty"`
	ReferredID    uint           `gorm:"not null" json:"referredId"`
	Referred      User           `gorm:"foreignKey:ReferredID" json:"referred,omitempty"`
	BusinessID    uint           `gorm:"not null" json:"businessId"`
	Business      Business       `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	BonusProgramID uint          `gorm:"not null" json:"bonusProgramId"`
	BonusProgram  BonusProgram  `gorm:"foreignKey:BonusProgramID" json:"bonusProgram,omitempty"`
	ReferralCode  string         `gorm:"not null" json:"referralCode"`
	Status        ReferralStatus `gorm:"not null;default:'pending';type:varchar(20)" json:"status"`
	ReferrerRewarded bool        `gorm:"default:false" json:"referrerRewarded"`
	ReferredRewarded bool        `gorm:"default:false" json:"referredRewarded"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

func (Referral) TableName() string {
	return "referrals"
}
