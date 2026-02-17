package postgres

import (
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type cityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) repository.CityRepository {
	return &cityRepository{db: db}
}

func (r *cityRepository) GetAll() ([]entity.City, error) {
	var cities []entity.City
	err := r.db.Order("name ASC").Find(&cities).Error
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func (r *cityRepository) GetByID(id uint) (*entity.City, error) {
	var city entity.City
	err := r.db.First(&city, id).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *cityRepository) GetBySlug(slug string) (*entity.City, error) {
	var city entity.City
	err := r.db.Where("slug = ?", slug).First(&city).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *cityRepository) Create(city *entity.City) error {
	return r.db.Create(city).Error
}

func (r *cityRepository) Update(city *entity.City) error {
	return r.db.Save(city).Error
}

func (r *cityRepository) Delete(id uint) error {
	return r.db.Delete(&entity.City{}, id).Error
}
