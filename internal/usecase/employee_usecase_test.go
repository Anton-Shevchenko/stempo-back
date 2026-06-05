package usecase

import (
	"strings"
	"testing"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) Create(employee *entity.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) FindByID(id uint) (*entity.Employee, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) FindByBusinessID(businessID uint) ([]entity.Employee, error) {
	args := m.Called(businessID)
	return args.Get(0).([]entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) FindByUserID(userID uint) ([]entity.Employee, error) {
	args := m.Called(userID)
	return args.Get(0).([]entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) FindByBusinessAndUser(businessID, userID uint) (*entity.Employee, error) {
	args := m.Called(businessID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmployeeRepository) DeleteByBusinessAndUser(businessID, userID uint) error {
	args := m.Called(businessID, userID)
	return args.Error(0)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendInviteEmail(toEmail, toName, inviteLink string) error {
	args := m.Called(toEmail, toName, inviteLink)
	return args.Error(0)
}

func TestEmployeeUsecase_AddEmployee_InvitesNewUser(t *testing.T) {
	const (
		businessID    uint = 1
		ownerID       uint = 5
		employeeEmail      = "new.employee@example.com"
	)

	mockBusinessRepo := new(MockBusinessRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	mockUserRepo := new(MockUserRepository)
	mockEmailSvc := new(MockEmailService)

	mockBusinessRepo.On("FindByID", businessID).Return(&entity.Business{
		ID:      businessID,
		OwnerID: ownerID,
	}, nil)

	mockUserRepo.On("FindByEmail", employeeEmail).Return(nil, assert.AnError)

	mockUserRepo.On("Create", mock.AnythingOfType("*entity.User")).Run(func(args mock.Arguments) {
		user := args.Get(0).(*entity.User)
		assert.Equal(t, employeeEmail, user.Email)
		assert.NotEmpty(t, user.Password)
		assert.NotNil(t, user.InviteToken)
		assert.NotEmpty(t, *user.InviteToken)
		assert.NotNil(t, user.InviteTokenExpiry)
		user.ID = 42
	}).Return(nil)

	mockEmployeeRepo.On("Create", mock.AnythingOfType("*entity.Employee")).Run(func(args mock.Arguments) {
		employee := args.Get(0).(*entity.Employee)
		assert.Equal(t, businessID, employee.BusinessID)
		assert.Equal(t, uint(42), employee.UserID)
		employee.ID = 7
	}).Return(nil)

	expectedEmployee := &entity.Employee{
		ID:         7,
		BusinessID: businessID,
		UserID:     42,
	}
	mockEmployeeRepo.On("FindByID", uint(7)).Return(expectedEmployee, nil)

	mockEmailSvc.On(
		"SendInviteEmail",
		employeeEmail,
		"",
		mock.MatchedBy(func(link string) bool {
			return strings.HasPrefix(link, "stempo://set-password?token=") &&
				len(strings.TrimPrefix(link, "stempo://set-password?token=")) > 0
		}),
	).Return(nil)

	usecase := NewEmployeeUsecase(mockEmployeeRepo, mockBusinessRepo, mockUserRepo, mockEmailSvc)

	employee, err := usecase.AddEmployee(businessID, ownerID, employeeEmail)

	assert.NoError(t, err)
	assert.NotNil(t, employee)
	assert.Equal(t, expectedEmployee, employee)

	mockBusinessRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
	mockEmailSvc.AssertExpectations(t)
}
