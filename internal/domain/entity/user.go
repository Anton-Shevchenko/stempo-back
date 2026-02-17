package entity

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Name      *string   `json:"name,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	AvatarURL *string   `json:"avatarUrl,omitempty"`
	CityID    *uint     `json:"cityId,omitempty"`
	City      *City     `gorm:"foreignKey:CityID" json:"city,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}
