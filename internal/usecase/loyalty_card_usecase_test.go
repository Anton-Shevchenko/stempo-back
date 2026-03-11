package usecase

import (
	"testing"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLoyaltyCardRepository struct {
	mock.Mock
}

func (m *MockLoyaltyCardRepository) Create(card *entity.LoyaltyCard) error {
	args := m.Called(card)
	return args.Error(0)
}

func (m *MockLoyaltyCardRepository) FindByID(id uint) (*entity.LoyaltyCard, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LoyaltyCard), args.Error(1)
}

func (m *MockLoyaltyCardRepository) FindByUserID(userID uint) ([]entity.LoyaltyCard, error) {
	args := m.Called(userID)
	return args.Get(0).([]entity.LoyaltyCard), args.Error(1)
}

func (m *MockLoyaltyCardRepository) FindByBusinessID(businessID uint) ([]entity.LoyaltyCard, error) {
	args := m.Called(businessID)
	return args.Get(0).([]entity.LoyaltyCard), args.Error(1)
}

func (m *MockLoyaltyCardRepository) Update(card *entity.LoyaltyCard) error {
	args := m.Called(card)
	return args.Error(0)
}

func (m *MockLoyaltyCardRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLoyaltyCardRepository) FindByUserAndBusiness(userID, businessID uint) (*entity.LoyaltyCard, error) {
	args := m.Called(userID, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LoyaltyCard), args.Error(1)
}

type MockEmployeeUsecase struct {
	mock.Mock
}

func (m *MockEmployeeUsecase) AddEmployee(businessID, ownerID uint, employeeEmail string) (*entity.Employee, error) {
	args := m.Called(businessID, ownerID, employeeEmail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Employee), args.Error(1)
}

func (m *MockEmployeeUsecase) RemoveEmployee(businessID, employeeID, ownerID uint) error {
	args := m.Called(businessID, employeeID, ownerID)
	return args.Error(0)
}

func (m *MockEmployeeUsecase) GetEmployeesByBusiness(businessID, ownerID uint) ([]entity.Employee, error) {
	args := m.Called(businessID, ownerID)
	return args.Get(0).([]entity.Employee), args.Error(1)
}

func (m *MockEmployeeUsecase) IsEmployee(businessID, userID uint) (bool, error) {
	args := m.Called(businessID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmployeeUsecase) GetBusinessesByEmployee(userID uint) ([]entity.Business, error) {
	args := m.Called(userID)
	return args.Get(0).([]entity.Business), args.Error(1)
}

func TestLoyaltyCardUsecase_Create(t *testing.T) {
	tests := []struct {
		name        string
		card        *entity.LoyaltyCard
		mockSetup   func(*MockLoyaltyCardRepository)
		expectError bool
	}{
		{
			name: "successful create",
			card: &entity.LoyaltyCard{
				UserID:         1,
				BusinessID:     1,
				Stamps:         0,
				StampsRequired: 10,
			},
			mockSetup: func(m *MockLoyaltyCardRepository) {
				m.On("FindByUserAndBusiness", uint(1), uint(1)).Return(nil, assert.AnError)
				m.On("Create", mock.AnythingOfType("*entity.LoyaltyCard")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "card already exists",
			card: &entity.LoyaltyCard{
				UserID:     1,
				BusinessID: 1,
			},
			mockSetup: func(m *MockLoyaltyCardRepository) {
				existingCard := &entity.LoyaltyCard{ID: 1}
				m.On("FindByUserAndBusiness", uint(1), uint(1)).Return(existingCard, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockLoyaltyCardRepository)
			mockBusinessRepo := new(MockBusinessRepository)
			mockEmployeeUsecase := new(MockEmployeeUsecase)
			tt.mockSetup(mockRepo)

			usecase := NewLoyaltyCardUsecase(mockRepo, mockBusinessRepo, mockEmployeeUsecase)
			err := usecase.Create(tt.card)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLoyaltyCardUsecase_AddStamp(t *testing.T) {
	tests := []struct {
		name        string
		cardID      uint
		userID      uint
		cardUserID  uint
		expectError bool
	}{
		{
			name:        "successful add stamp",
			cardID:      1,
			userID:      1,
			cardUserID:  1,
			expectError: false,
		},
		{
			name:        "unauthorized add stamp",
			cardID:      1,
			userID:      2,
			cardUserID:  1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockLoyaltyCardRepository)
			mockBusinessRepo := new(MockBusinessRepository)
			mockEmployeeUsecase := new(MockEmployeeUsecase)

			card := &entity.LoyaltyCard{
				ID:             tt.cardID,
				UserID:         tt.cardUserID,
				Stamps:         5,
				StampsRequired: 10,
				BusinessID:     1,
			}

			mockRepo.On("FindByID", tt.cardID).Return(card, nil)
			// Business lookup happens only when userID != card.UserID (employee/owner case),
			// which we don't simulate here beyond unauthorized error path.
			mockBusiness := &entity.Business{ID: 1, OwnerID: 999}
			mockBusinessRepo.On("FindByID", uint(1)).Return(mockBusiness, nil)

			// For unauthorized case, IsEmployee should return false without error
			if tt.expectError {
				mockEmployeeUsecase.On("IsEmployee", uint(1), tt.userID).Return(false, nil)
			}

			if !tt.expectError {
				mockRepo.On("Update", mock.AnythingOfType("*entity.LoyaltyCard")).Return(nil)
			}

			usecase := NewLoyaltyCardUsecase(mockRepo, mockBusinessRepo, mockEmployeeUsecase)
			err := usecase.AddStamp(tt.cardID, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLoyaltyCardUsecase_GetByUserID(t *testing.T) {
	mockRepo := new(MockLoyaltyCardRepository)
	mockBusinessRepo := new(MockBusinessRepository)
	mockEmployeeUsecase := new(MockEmployeeUsecase)
	cards := []entity.LoyaltyCard{
		{ID: 1, UserID: 1, BusinessID: 1},
		{ID: 2, UserID: 1, BusinessID: 2},
	}

	mockRepo.On("FindByUserID", uint(1)).Return(cards, nil)

	usecase := NewLoyaltyCardUsecase(mockRepo, mockBusinessRepo, mockEmployeeUsecase)
	result, err := usecase.GetByUserID(1)

	assert.NoError(t, err)
	assert.Equal(t, cards, result)
	mockRepo.AssertExpectations(t)
}
