package usecase

import (
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
)

type CategoryUsecase interface {
	GetAll() ([]entity.Category, error)
	GetBySlug(slug string) (*entity.Category, error)
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{categoryRepo: categoryRepo}
}

func (u *categoryUsecase) GetAll() ([]entity.Category, error) {
	return u.categoryRepo.GetAll()
}

func (u *categoryUsecase) GetBySlug(slug string) (*entity.Category, error) {
	return u.categoryRepo.GetBySlug(slug)
}
