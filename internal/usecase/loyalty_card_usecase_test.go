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

func (m *MockLoyaltyCardRepository) FindByUserAndBusiness(userID, businessID uint) (*entity.LoyaltyCard, error) {
	args := m.Called(userID, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LoyaltyCard), args.Error(1)
}

func (m *MockLoyaltyCardRepository) Update(card *entity.LoyaltyCard) error {
	args := m.Called(card)
	return args.Error(0)
}

func (m *MockLoyaltyCardRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
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
				UserID:         1,
				BusinessID:     1,
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
			tt.mockSetup(mockRepo)

			usecase := NewLoyaltyCardUsecase(mockRepo)
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

			card := &entity.LoyaltyCard{
				ID:             tt.cardID,
				UserID:         tt.cardUserID,
				Stamps:         5,
				StampsRequired: 10,
			}

			mockRepo.On("FindByID", tt.cardID).Return(card, nil)
			if !tt.expectError {
				mockRepo.On("Update", mock.AnythingOfType("*entity.LoyaltyCard")).Return(nil)
			}

			usecase := NewLoyaltyCardUsecase(mockRepo)
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
	cards := []entity.LoyaltyCard{
		{ID: 1, UserID: 1, BusinessID: 1},
		{ID: 2, UserID: 1, BusinessID: 2},
	}

	mockRepo.On("FindByUserID", uint(1)).Return(cards, nil)

	usecase := NewLoyaltyCardUsecase(mockRepo)
	result, err := usecase.GetByUserID(1)

	assert.NoError(t, err)
	assert.Equal(t, cards, result)
	mockRepo.AssertExpectations(t)
}
