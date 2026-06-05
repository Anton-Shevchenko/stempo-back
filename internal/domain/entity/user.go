package entity

import "time"

type User struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	Email               string     `gorm:"uniqueIndex;not null" json:"email"`
	Password            string     `gorm:"default:''" json:"-"`
	GoogleID            *string    `gorm:"uniqueIndex" json:"-"`
	Name                *string    `json:"name,omitempty"`
	Phone               *string    `json:"phone,omitempty"`
	AvatarURL           *string    `json:"avatarUrl,omitempty"`
	CityID              *uint      `json:"cityId,omitempty"`
	City                *City      `gorm:"foreignKey:CityID" json:"city,omitempty"`
	InviteToken         *string    `gorm:"uniqueIndex" json:"-"`
	InviteTokenExpiry   *time.Time `json:"-"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}
