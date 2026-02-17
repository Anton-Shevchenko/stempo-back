package repository

import "github.com/stempo/backend/internal/domain/entity"

type CategoryRepository interface {
	GetAll() ([]entity.Category, error)
	GetBySlug(slug string) (*entity.Category, error)
	Create(category *entity.Category) error
}
