package entity

import "time"

type BusinessStatus string

const (
	BusinessStatusPending  BusinessStatus = "pending"
	BusinessStatusApproved BusinessStatus = "approved"
	BusinessStatusRejected BusinessStatus = "rejected"
)

type Business struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	Name            string        `gorm:"not null" json:"name"`
	Category        string        `gorm:"not null" json:"category"`
	Address         string        `gorm:"not null" json:"address"`
	Rating          float64       `gorm:"default:0" json:"rating"`
	IsOpen          bool          `gorm:"default:true" json:"isOpen"`
	ImageURL        *string       `json:"imageUrl,omitempty"`
	Icon            string        `json:"icon"`
	IconColor       string        `json:"iconColor"`
	Description     *string       `json:"description,omitempty"`
	HasLoyaltyProgram bool        `gorm:"default:false" json:"hasLoyaltyProgram"`
	Featured        bool          `gorm:"default:false" json:"featured"`
	Status          BusinessStatus `gorm:"default:'pending'" json:"status"`
	RejectionReason *string       `json:"rejectionReason,omitempty"`
	OwnerID         uint          `gorm:"not null" json:"ownerId"`
	Owner           User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

func (Business) TableName() string {
	return "businesses"
}
