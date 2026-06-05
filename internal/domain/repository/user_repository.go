package repository

import "github.com/stempo/backend/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uint) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByGoogleID(googleID string) (*entity.User, error)
	FindByInviteToken(token string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
}
