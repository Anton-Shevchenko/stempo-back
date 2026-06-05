package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/domain/repository"
	"github.com/stempo/backend/internal/infrastructure/email"
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
	emailSvc     email.EmailService
}

func NewEmployeeUsecase(
	employeeRepo repository.EmployeeRepository,
	businessRepo repository.BusinessRepository,
	userRepo repository.UserRepository,
	emailSvc email.EmailService,
) EmployeeUsecase {
	return &employeeUsecase{
		employeeRepo: employeeRepo,
		businessRepo: businessRepo,
		userRepo:     userRepo,
		emailSvc:     emailSvc,
	}
}

func (u *employeeUsecase) AddEmployee(businessID, ownerID uint, employeeEmail string) (*entity.Employee, error) {
	business, err := u.businessRepo.FindByID(businessID)
	if err != nil {
		return nil, errors.New("business not found")
	}

	if business.OwnerID != ownerID {
		return nil, errors.New("unauthorized: only business owner can add employees")
	}

	employeeUser, err := u.userRepo.FindByEmail(employeeEmail)
	if err != nil || employeeUser == nil {
		// User does not exist — create a pending invited user and send invite email.
		invitedEmployee, inviteErr := u.createInvitedEmployee(businessID, ownerID, employeeEmail)
		if inviteErr != nil {
			return nil, inviteErr
		}
		return invitedEmployee, nil
	}

	// Check if user is already an employee
	existing, _ := u.employeeRepo.FindByBusinessAndUser(businessID, employeeUser.ID)
	if existing != nil {
		return nil, errors.New("user is already an employee")
	}

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

	return u.employeeRepo.FindByID(employee.ID)
}

func (u *employeeUsecase) createInvitedEmployee(businessID, ownerID uint, employeeEmail string) (*entity.Employee, error) {
	fmt.Println("createInvitedEmployee")
	token, err := generateSecureToken()
	if err != nil {
		fmt.Println("failed to generate invite token: %w", err)
		return nil, fmt.Errorf("failed to generate invite token: %w", err)
	}

	// Use a random bcrypt hash as placeholder so the not-null constraint is satisfied.
	// The user cannot log in until they set a real password via the invite link.
	randomBytes := make([]byte, 16)
	if _, randErr := rand.Read(randomBytes); randErr != nil {
		return nil, fmt.Errorf("failed to generate placeholder password: %w", randErr)
	}
	placeholderPassword, err := bcrypt.GenerateFromPassword(randomBytes, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash placeholder password: %w", err)
	}

	expiry := time.Now().Add(48 * time.Hour)
	newUser := &entity.User{
		Email:               employeeEmail,
		Password:            string(placeholderPassword),
		InviteToken:         &token,
		InviteTokenExpiry:   &expiry,
	}

	if err := u.userRepo.Create(newUser); err != nil {
		return nil, fmt.Errorf("failed to create invited user: %w", err)
	}

	employee := &entity.Employee{
		BusinessID: businessID,
		UserID:     newUser.ID,
	}

	if err := u.employeeRepo.Create(employee); err != nil {
		return nil, fmt.Errorf("failed to create employee record: %w", err)
	}

	deepLinkBase := os.Getenv("APP_DEEP_LINK_BASE")
	if deepLinkBase == "" {
		deepLinkBase = "stempo://"
	}
	inviteLink := fmt.Sprintf("%sset-password?token=%s", deepLinkBase, token)

	if sendErr := u.emailSvc.SendInviteEmail(employeeEmail, "", inviteLink); sendErr != nil {
		log.Printf("warning: failed to send invite email to %s: %v", employeeEmail, sendErr)
	}

	return u.employeeRepo.FindByID(employee.ID)
}

func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (u *employeeUsecase) RemoveEmployee(businessID, employeeID, ownerID uint) error {
	business, err := u.businessRepo.FindByID(businessID)
	if err != nil {
		return errors.New("business not found")
	}

	if business.OwnerID != ownerID {
		return errors.New("unauthorized: only business owner can remove employees")
	}

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
	fmt.Println("GetEmployeesByBusiness")
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
		return false, nil
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
