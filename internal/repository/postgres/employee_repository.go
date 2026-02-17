package postgres

import (
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) repository.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(employee *entity.Employee) error {
	return r.db.Create(employee).Error
}

func (r *employeeRepository) FindByID(id uint) (*entity.Employee, error) {
	var employee entity.Employee
	err := r.db.Preload("Business").Preload("User").First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) FindByBusinessID(businessID uint) ([]entity.Employee, error) {
	var employees []entity.Employee
	err := r.db.Where("business_id = ?", businessID).
		Preload("User").
		Preload("Business").
		Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) FindByUserID(userID uint) ([]entity.Employee, error) {
	var employees []entity.Employee
	err := r.db.Where("user_id = ?", userID).
		Preload("Business").
		Preload("User").
		Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) FindByBusinessAndUser(businessID, userID uint) (*entity.Employee, error) {
	var employee entity.Employee
	err := r.db.Where("business_id = ? AND user_id = ?", businessID, userID).
		Preload("Business").
		Preload("User").
		First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Employee{}, id).Error
}

func (r *employeeRepository) DeleteByBusinessAndUser(businessID, userID uint) error {
	return r.db.Where("business_id = ? AND user_id = ?", businessID, userID).
		Delete(&entity.Employee{}).Error
}
