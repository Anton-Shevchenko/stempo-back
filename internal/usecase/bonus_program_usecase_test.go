package usecase

import (
	"testing"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBonusProgramRepository struct {
	mock.Mock
}

func (m *MockBonusProgramRepository) Create(program *entity.BonusProgram) error {
	args := m.Called(program)
	return args.Error(0)
}

func (m *MockBonusProgramRepository) FindByID(id uint) (*entity.BonusProgram, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BonusProgram), args.Error(1)
}

func (m *MockBonusProgramRepository) FindByBusinessID(businessID uint) ([]entity.BonusProgram, error) {
	args := m.Called(businessID)
	return args.Get(0).([]entity.BonusProgram), args.Error(1)
}

func (m *MockBonusProgramRepository) FindAll(page, pageSize int) ([]entity.BonusProgram, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]entity.BonusProgram), args.Get(1).(int64), args.Error(2)
}

func (m *MockBonusProgramRepository) FindByStatus(status entity.BonusProgramStatus, page, pageSize int) ([]entity.BonusProgram, int64, error) {
	args := m.Called(status, page, pageSize)
	return args.Get(0).([]entity.BonusProgram), args.Get(1).(int64), args.Error(2)
}

func (m *MockBonusProgramRepository) Update(program *entity.BonusProgram) error {
	args := m.Called(program)
	return args.Error(0)
}

func (m *MockBonusProgramRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBonusProgramUsecase_Create(t *testing.T) {
	tests := []struct {
		name        string
		businessID  uint
		userID      uint
		ownerID     uint
		expectError bool
	}{
		{
			name:        "successful create",
			businessID:  1,
			userID:      1,
			ownerID:     1,
			expectError: false,
		},
		{
			name:        "unauthorized create",
			businessID:  1,
			userID:      2,
			ownerID:     1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProgramRepo := new(MockBonusProgramRepository)
			mockBusinessRepo := new(MockBusinessRepository)

			program := &entity.BonusProgram{
				BusinessID:     tt.businessID,
				Name:           "Test Program",
				PointsRequired: 10,
			}

			business := &entity.Business{
				ID:      tt.businessID,
				OwnerID: tt.ownerID,
			}

			mockBusinessRepo.On("FindByID", tt.businessID).Return(business, nil)
			if !tt.expectError {
				mockProgramRepo.On("Create", program).Return(nil)
			}

			usecase := NewBonusProgramUsecase(mockProgramRepo, mockBusinessRepo)
			err := usecase.Create(program, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockProgramRepo.AssertExpectations(t)
			mockBusinessRepo.AssertExpectations(t)
		})
	}
}

func TestBonusProgramUsecase_GetByBusinessID(t *testing.T) {
	mockProgramRepo := new(MockBonusProgramRepository)
	mockBusinessRepo := new(MockBusinessRepository)

	programs := []entity.BonusProgram{
		{ID: 1, BusinessID: 1, Name: "Program 1"},
		{ID: 2, BusinessID: 1, Name: "Program 2"},
	}

	mockProgramRepo.On("FindByBusinessID", uint(1)).Return(programs, nil)

	usecase := NewBonusProgramUsecase(mockProgramRepo, mockBusinessRepo)
	result, err := usecase.GetByBusinessID(1)

	assert.NoError(t, err)
	assert.Equal(t, programs, result)
	mockProgramRepo.AssertExpectations(t)
}

func TestBonusProgramUsecase_Delete(t *testing.T) {
	tests := []struct {
		name        string
		programID   uint
		userID      uint
		businessID  uint
		ownerID     uint
		expectError bool
	}{
		{
			name:        "successful delete",
			programID:   1,
			userID:      1,
			businessID:  1,
			ownerID:     1,
			expectError: false,
		},
		{
			name:        "unauthorized delete",
			programID:   1,
			userID:      2,
			businessID:  1,
			ownerID:     1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProgramRepo := new(MockBonusProgramRepository)
			mockBusinessRepo := new(MockBusinessRepository)

			program := &entity.BonusProgram{
				ID:         tt.programID,
				BusinessID: tt.businessID,
			}

			business := &entity.Business{
				ID:      tt.businessID,
				OwnerID: tt.ownerID,
			}

			mockProgramRepo.On("FindByID", tt.programID).Return(program, nil)
			mockBusinessRepo.On("FindByID", tt.businessID).Return(business, nil)
			if !tt.expectError {
				mockProgramRepo.On("Delete", tt.programID).Return(nil)
			}

			usecase := NewBonusProgramUsecase(mockProgramRepo, mockBusinessRepo)
			err := usecase.Delete(tt.programID, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockProgramRepo.AssertExpectations(t)
			mockBusinessRepo.AssertExpectations(t)
		})
	}
}
