package entity

import "time"

type BonusProgramStatus string

const (
	BonusProgramStatusPending  BonusProgramStatus = "pending"
	BonusProgramStatusApproved BonusProgramStatus = "approved"
	BonusProgramStatusRejected BonusProgramStatus = "rejected"
	BonusProgramStatusActive   BonusProgramStatus = "active"
)

type QRCodeType string

const (
	QRCodeTypeTemporary QRCodeType = "temporary"
	QRCodeTypePermanent QRCodeType = "permanent"
)

type ProgramType string

const (
	ProgramTypeStampCard    ProgramType = "stamp_card"   // Punch/Stamp Card - collect stamps for rewards
	ProgramTypePoints       ProgramType = "points"       // Points-Based - earn points per purchase
	ProgramTypeTiered       ProgramType = "tiered"       // Tiered - different levels with increasing benefits
	ProgramTypeReferral     ProgramType = "referral"     // Referral - rewards for referring friends
	ProgramTypeGamification ProgramType = "gamification" // Gamification - challenges and achievements
	ProgramTypeVIP          ProgramType = "vip"          // VIP/Member Benefits - exclusive perks
	ProgramTypeSpendBased   ProgramType = "spend_based"  // Spend-Based - rewards based on spending amount
)

type BonusProgram struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Name        string      `gorm:"not null" json:"name"`
	Description string      `gorm:"not null" json:"description"`
	BusinessID  uint        `gorm:"not null" json:"businessId"`
	Business    Business    `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	ProgramType ProgramType `gorm:"not null;default:'stamp_card';type:varchar(20)" json:"programType"`

	// Stamp/Points configuration
	Points            int      `gorm:"default:0" json:"points"`
	PointsRequired    int      `gorm:"not null" json:"pointsRequired"`
	PointsPerPurchase int      `gorm:"default:1" json:"pointsPerPurchase"` // For points-based programs
	PointsPerDollar   *float64 `json:"pointsPerDollar,omitempty"`          // Points earned per dollar spent

	// Reward configuration
	Discount     int    `gorm:"not null" json:"discount"`
	DiscountType string `gorm:"not null;type:varchar(20)" json:"discountType"`

	// Tiered program configuration
	TierLevels     *string `json:"tierLevels,omitempty"`     // JSON array of tier configurations
	TierThresholds *string `json:"tierThresholds,omitempty"` // JSON array of spending thresholds

	// Referral program configuration
	ReferralReward *int `json:"referralReward,omitempty"` // Reward for referrer
	ReferredReward *int `json:"referredReward,omitempty"` // Reward for new customer

	// Gamification configuration
	Challenges   *string `json:"challenges,omitempty"`   // JSON array of challenges
	Achievements *string `json:"achievements,omitempty"` // JSON array of achievements

	// VIP/Member benefits
	MemberBenefits *string `json:"memberBenefits,omitempty"` // JSON array of benefits

	Status              BonusProgramStatus `gorm:"not null;default:'pending';type:varchar(20)" json:"status"`
	QRCodeType          QRCodeType         `gorm:"not null;default:'temporary';type:varchar(20)" json:"qrCodeType"`
	QRExpirationMinutes int                `gorm:"not null;default:3" json:"qrExpirationMinutes"` // 1-10 min for temporary QR
	RejectionReason     *string            `json:"rejectionReason,omitempty"`
	ExpiresAt           *time.Time         `json:"expiresAt,omitempty"`
	ImageURL            *string            `json:"imageUrl,omitempty"`
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`
}

func (BonusProgram) TableName() string {
	return "bonus_programs"
}
