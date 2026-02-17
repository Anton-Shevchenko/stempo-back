package usecase

import (
	"os"
	"testing"

	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAuthUsecase_Register(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name          string
		email         string
		password      string
		namePtr       *string
		phonePtr      *string
		mockSetup     func(*MockUserRepository)
		expectError   bool
		expectUser    bool
		expectToken   bool
	}{
		{
			name:     "successful registration",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(m *MockUserRepository) {
				m.On("FindByEmail", "test@example.com").Return(nil, assert.AnError)
				m.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)
			},
			expectError: false,
			expectUser:  true,
			expectToken: true,
		},
		{
			name:     "user already exists",
			email:    "existing@example.com",
			password: "password123",
			mockSetup: func(m *MockUserRepository) {
				existingUser := &entity.User{Email: "existing@example.com"}
				m.On("FindByEmail", "existing@example.com").Return(existingUser, nil)
			},
			expectError: true,
			expectUser:  false,
			expectToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			usecase := NewAuthUsecase(mockRepo)
			user, token, refreshToken, err := usecase.Register(tt.email, tt.password, tt.namePtr, tt.phonePtr, nil)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Empty(t, token)
				assert.Empty(t, refreshToken)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.NotEmpty(t, token)
				assert.NotEmpty(t, refreshToken)
				assert.Equal(t, tt.email, user.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUsecase_Login(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name        string
		email       string
		password    string
		mockSetup   func(*MockUserRepository)
		expectError bool
		expectUser  bool
		expectToken bool
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(m *MockUserRepository) {
				hashedPassword, _ := hashPassword("password123")
				user := &entity.User{
					ID:       1,
					Email:    "test@example.com",
					Password: hashedPassword,
				}
				m.On("FindByEmail", "test@example.com").Return(user, nil)
			},
			expectError: false,
			expectUser:  true,
			expectToken: true,
		},
		{
			name:     "user not found",
			email:    "notfound@example.com",
			password: "password123",
			mockSetup: func(m *MockUserRepository) {
				m.On("FindByEmail", "notfound@example.com").Return(nil, assert.AnError)
			},
			expectError: true,
			expectUser:  false,
			expectToken: false,
		},
		{
			name:     "wrong password",
			email:    "test@example.com",
			password: "wrongpassword",
			mockSetup: func(m *MockUserRepository) {
				hashedPassword, _ := hashPassword("password123")
				user := &entity.User{
					ID:       1,
					Email:    "test@example.com",
					Password: hashedPassword,
				}
				m.On("FindByEmail", "test@example.com").Return(user, nil)
			},
			expectError: true,
			expectUser:  false,
			expectToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			usecase := NewAuthUsecase(mockRepo)
			user, token, refreshToken, err := usecase.Login(tt.email, tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Empty(t, token)
				assert.Empty(t, refreshToken)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.NotEmpty(t, token)
				assert.NotEmpty(t, refreshToken)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUsecase_GetCurrentUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	user := &entity.User{
		ID:    1,
		Email: "test@example.com",
	}
	mockRepo.On("FindByID", uint(1)).Return(user, nil)

	usecase := NewAuthUsecase(mockRepo)
	result, err := usecase.GetCurrentUser(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Email, result.Email)
	mockRepo.AssertExpectations(t)
}
