package usecase

import (
	"errors"
	"strings"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
)

type CityUsecase interface {
	GetAll() ([]entity.City, error)
	GetByID(id uint) (*entity.City, error)
	Create(name string, slug string, latitude, longitude *float64) (*entity.City, error)
	Update(id uint, name string, slug string, latitude, longitude *float64) (*entity.City, error)
	Delete(id uint) error
}

type cityUsecase struct {
	cityRepo repository.CityRepository
}

func NewCityUsecase(cityRepo repository.CityRepository) CityUsecase {
	return &cityUsecase{cityRepo: cityRepo}
}

func (u *cityUsecase) GetAll() ([]entity.City, error) {
	return u.cityRepo.GetAll()
}

func (u *cityUsecase) GetByID(id uint) (*entity.City, error) {
	return u.cityRepo.GetByID(id)
}

func (u *cityUsecase) Create(name string, slug string, latitude, longitude *float64) (*entity.City, error) {
	if name == "" {
		return nil, errors.New("city name is required")
	}

	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	}

	existingCity, _ := u.cityRepo.GetBySlug(slug)
	if existingCity != nil {
		return nil, errors.New("city with this slug already exists")
	}

	city := &entity.City{
		Name:      name,
		Slug:      slug,
		Latitude:  latitude,
		Longitude: longitude,
	}

	if err := u.cityRepo.Create(city); err != nil {
		return nil, err
	}

	return city, nil
}

func (u *cityUsecase) Update(id uint, name string, slug string, latitude, longitude *float64) (*entity.City, error) {
	city, err := u.cityRepo.GetByID(id)
	if err != nil || city == nil {
		return nil, errors.New("city not found")
	}

	if name != "" {
		city.Name = name
	}

	if slug != "" {
		existingCity, _ := u.cityRepo.GetBySlug(slug)
		if existingCity != nil && existingCity.ID != id {
			return nil, errors.New("city with this slug already exists")
		}
		city.Slug = slug
	}

	if latitude != nil {
		city.Latitude = latitude
	}

	if longitude != nil {
		city.Longitude = longitude
	}

	if err := u.cityRepo.Update(city); err != nil {
		return nil, err
	}

	return city, nil
}

func (u *cityUsecase) Delete(id uint) error {
	city, err := u.cityRepo.GetByID(id)
	if err != nil || city == nil {
		return errors.New("city not found")
	}

	return u.cityRepo.Delete(id)
}
