package usecase

import (
	"errors"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
)

type EmployeeUsecase interface {
	AddEmployee(businessID, ownerID uint, employeeEmail string) (*entity.Employee, error)
	RemoveEmployee(businessID, employeeID, ownerID uint) error
	GetEmployeesByBusiness(businessID, ownerID uint) ([]entity.Employee, error)
	IsEmployee(businessID, userID uint) (bool, error)
	GetBusinessesByEmployee(userID uint) ([]entity.Business, error)
}

type employeeUsecase struct {
	employeeRepo repository.EmployeeRepository
	businessRepo repository.BusinessRepository
	userRepo     repository.UserRepository
}

func NewEmployeeUsecase(
	employeeRepo repository.EmployeeRepository,
	businessRepo repository.BusinessRepository,
	userRepo repository.UserRepository,
) EmployeeUsecase {
	return &employeeUsecase{
		employeeRepo: employeeRepo,
		businessRepo: businessRepo,
		userRepo:     userRepo,
	}
}

func (u *employeeUsecase) AddEmployee(businessID, ownerID uint, employeeEmail string) (*entity.Employee, error) {
	// Verify that the requester is the owner of the business
	business, err := u.businessRepo.FindByID(businessID)
	if err != nil {
		return nil, errors.New("business not found")
	}

	if business.OwnerID != ownerID {
		return nil, errors.New("unauthorized: only business owner can add employees")
	}

	// Find user by email
	employeeUser, err := u.userRepo.FindByEmail(employeeEmail)
	if err != nil || employeeUser == nil {
		return nil, errors.New("user not found")
	}

	// Check if user is already an employee
	existing, _ := u.employeeRepo.FindByBusinessAndUser(businessID, employeeUser.ID)
	if existing != nil {
		return nil, errors.New("user is already an employee")
	}

	// Prevent adding owner as employee
	if employeeUser.ID == ownerID {
		return nil, errors.New("business owner cannot be added as employee")
	}

	employee := &entity.Employee{
		BusinessID: businessID,
		UserID:     employeeUser.ID,
	}

	if err := u.employeeRepo.Create(employee); err != nil {
		return nil, err
	}

	// Reload with relations
	return u.employeeRepo.FindByID(employee.ID)
}

func (u *employeeUsecase) RemoveEmployee(businessID, employeeID, ownerID uint) error {
	// Verify that the requester is the owner of the business
	business, err := u.businessRepo.FindByID(businessID)
	if err != nil {
		return errors.New("business not found")
	}

	if business.OwnerID != ownerID {
		return errors.New("unauthorized: only business owner can remove employees")
	}

	// Verify employee exists and belongs to this business
	employee, err := u.employeeRepo.FindByID(employeeID)
	if err != nil {
		return errors.New("employee not found")
	}

	if employee.BusinessID != businessID {
		return errors.New("employee does not belong to this business")
	}

	return u.employeeRepo.Delete(employeeID)
}

func (u *employeeUsecase) GetEmployeesByBusiness(businessID, ownerID uint) ([]entity.Employee, error) {
	// Verify that the requester is the owner of the business
	business, err := u.businessRepo.FindByID(businessID)
	if err != nil {
		return nil, errors.New("business not found")
	}

	if business.OwnerID != ownerID {
		return nil, errors.New("unauthorized: only business owner can view employees")
	}

	return u.employeeRepo.FindByBusinessID(businessID)
}

func (u *employeeUsecase) IsEmployee(businessID, userID uint) (bool, error) {
	employee, err := u.employeeRepo.FindByBusinessAndUser(businessID, userID)
	if err != nil {
		return false, nil // Not an error, just not an employee
	}
	return employee != nil, nil
}

func (u *employeeUsecase) GetBusinessesByEmployee(userID uint) ([]entity.Business, error) {
	employees, err := u.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	businesses := make([]entity.Business, 0, len(employees))
	for _, emp := range employees {
		business, err := u.businessRepo.FindByID(emp.BusinessID)
		if err == nil && business != nil {
			businesses = append(businesses, *business)
		}
	}

	return businesses, nil
}
