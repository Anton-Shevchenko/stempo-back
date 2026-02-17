package usecase

import (
	"testing"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBusinessRepository struct {
	mock.Mock
}

func (m *MockBusinessRepository) Create(business *entity.Business) error {
	args := m.Called(business)
	return args.Error(0)
}

func (m *MockBusinessRepository) FindByID(id uint) (*entity.Business, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Business), args.Error(1)
}

func (m *MockBusinessRepository) FindByOwnerID(ownerID uint) (*entity.Business, error) {
	args := m.Called(ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Business), args.Error(1)
}

func (m *MockBusinessRepository) FindFeatured() ([]entity.Business, error) {
	args := m.Called()
	return args.Get(0).([]entity.Business), args.Error(1)
}

func (m *MockBusinessRepository) FindAll(page, pageSize int) ([]entity.Business, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]entity.Business), args.Get(1).(int64), args.Error(2)
}

func (m *MockBusinessRepository) Update(business *entity.Business) error {
	args := m.Called(business)
	return args.Error(0)
}

func (m *MockBusinessRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBusinessUsecase_Create(t *testing.T) {
	mockRepo := new(MockBusinessRepository)
	business := &entity.Business{
		Name:     "Test Business",
		Category: "coffee",
		Address:  "123 Main St",
		OwnerID:  1,
	}

	mockRepo.On("Create", business).Return(nil)

	usecase := NewBusinessUsecase(mockRepo)
	err := usecase.Create(business)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBusinessUsecase_Update(t *testing.T) {
	tests := []struct {
		name        string
		businessID  uint
		userID      uint
		ownerID     uint
		expectError bool
	}{
		{
			name:        "successful update",
			businessID:  1,
			userID:      1,
			ownerID:     1,
			expectError: false,
		},
		{
			name:        "unauthorized update",
			businessID:  1,
			userID:      2,
			ownerID:     1,
			expectError: true,
		},
		{
			name:        "business not found",
			businessID:  999,
			userID:      1,
			ownerID:     1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBusinessRepository)
			business := &entity.Business{
				ID:      tt.businessID,
				Name:    "Updated Business",
				OwnerID: tt.ownerID,
			}

			if tt.expectError && tt.businessID == 999 {
				mockRepo.On("FindByID", tt.businessID).Return(nil, assert.AnError)
			} else {
				existing := &entity.Business{
					ID:      tt.businessID,
					OwnerID: tt.ownerID,
				}
				mockRepo.On("FindByID", tt.businessID).Return(existing, nil)
				if !tt.expectError {
					mockRepo.On("Update", mock.AnythingOfType("*entity.Business")).Return(nil)
				}
			}

			usecase := NewBusinessUsecase(mockRepo)
			err := usecase.Update(business, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBusinessUsecase_Delete(t *testing.T) {
	tests := []struct {
		name        string
		businessID  uint
		userID      uint
		ownerID     uint
		expectError bool
	}{
		{
			name:        "successful delete",
			businessID:  1,
			userID:      1,
			ownerID:     1,
			expectError: false,
		},
		{
			name:        "unauthorized delete",
			businessID:  1,
			userID:      2,
			ownerID:     1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBusinessRepository)
			business := &entity.Business{
				ID:      tt.businessID,
				OwnerID: tt.ownerID,
			}

			mockRepo.On("FindByID", tt.businessID).Return(business, nil)
			if !tt.expectError {
				mockRepo.On("Delete", tt.businessID).Return(nil)
			}

			usecase := NewBusinessUsecase(mockRepo)
			err := usecase.Delete(tt.businessID, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBusinessUsecase_GetAll(t *testing.T) {
	mockRepo := new(MockBusinessRepository)
	businesses := []entity.Business{
		{ID: 1, Name: "Business 1"},
		{ID: 2, Name: "Business 2"},
	}
	total := int64(2)

	mockRepo.On("FindAll", 1, 10).Return(businesses, total, nil)

	usecase := NewBusinessUsecase(mockRepo)
	result, totalResult, err := usecase.GetAll(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, businesses, result)
	assert.Equal(t, total, totalResult)
	mockRepo.AssertExpectations(t)
}

func TestBusinessUsecase_GetAll_DefaultPagination(t *testing.T) {
	mockRepo := new(MockBusinessRepository)
	businesses := []entity.Business{}
	total := int64(0)

	mockRepo.On("FindAll", 1, 10).Return(businesses, total, nil)

	usecase := NewBusinessUsecase(mockRepo)
	result, totalResult, err := usecase.GetAll(0, 0)

	assert.NoError(t, err)
	assert.Equal(t, businesses, result)
	assert.Equal(t, total, totalResult)
	mockRepo.AssertExpectations(t)
}
