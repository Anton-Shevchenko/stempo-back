package repository

import "github.com/stempo/backend/internal/domain/entity"

type CityRepository interface {
	GetAll() ([]entity.City, error)
	GetByID(id uint) (*entity.City, error)
	GetBySlug(slug string) (*entity.City, error)
	Create(city *entity.City) error
	Update(city *entity.City) error
	Delete(id uint) error
}
