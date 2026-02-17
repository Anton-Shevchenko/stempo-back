package repository

import "github.com/stempo/backend/internal/domain/entity"

type EmployeeRepository interface {
	Create(employee *entity.Employee) error
	FindByID(id uint) (*entity.Employee, error)
	FindByBusinessID(businessID uint) ([]entity.Employee, error)
	FindByUserID(userID uint) ([]entity.Employee, error)
	FindByBusinessAndUser(businessID, userID uint) (*entity.Employee, error)
	Delete(id uint) error
	DeleteByBusinessAndUser(businessID, userID uint) error
}
