package entity

import "time"

type LoyaltyCard struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null" json:"userId"`
	User           User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BusinessID     uint      `gorm:"not null" json:"businessId"`
	Business       Business  `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	BonusProgramID *uint     `json:"bonusProgramId,omitempty"`
	BonusProgram   *BonusProgram `gorm:"foreignKey:BonusProgramID" json:"bonusProgram,omitempty"`
	
	// Stamp/Points tracking
	Stamps         int       `gorm:"default:0" json:"stamps"`
	StampsRequired int       `gorm:"not null" json:"stampsRequired"`
	Points         int       `gorm:"default:0" json:"points"` // For points-based programs
	TotalSpent     float64   `gorm:"default:0" json:"totalSpent"` // For spend-based programs
	
	// Tier tracking
	CurrentTier    *string   `json:"currentTier,omitempty"` // Current tier level
	
	// Referral tracking
	ReferralCode   *string   `gorm:"uniqueIndex" json:"referralCode,omitempty"` // Unique referral code for this user
	ReferralsCount int       `gorm:"default:0" json:"referralsCount"` // Number of successful referrals
	
	// Gamification tracking
	Achievements   *string   `json:"achievements,omitempty"` // JSON array of unlocked achievements
	ChallengesProgress *string `json:"challengesProgress,omitempty"` // JSON object tracking challenge progress
	
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (LoyaltyCard) TableName() string {
	return "loyalty_cards"
}
