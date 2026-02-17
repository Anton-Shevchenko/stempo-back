package entity

import "time"

type Employee struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	BusinessID uint      `gorm:"not null;index" json:"businessId"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	UserID     uint      `gorm:"not null;index" json:"userId"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (Employee) TableName() string {
	return "employees"
}
