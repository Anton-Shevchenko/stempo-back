package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Register(email, password string, name, phone *string, cityID *uint) (*entity.User, string, string, error) {
	args := m.Called(email, password, name, phone, cityID)
	if args.Get(0) == nil {
		return nil, "", "", args.Error(3)
	}
	return args.Get(0).(*entity.User), args.String(1), args.String(2), args.Error(3)
}

func (m *MockAuthUsecase) Login(email, password string) (*entity.User, string, string, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, "", "", args.Error(3)
	}
	return args.Get(0).(*entity.User), args.String(1), args.String(2), args.Error(3)
}

func (m *MockAuthUsecase) LoginWithGoogle(idToken string) (*entity.User, string, string, error) {
	args := m.Called(idToken)
	if args.Get(0) == nil {
		return nil, "", "", args.Error(3)
	}
	return args.Get(0).(*entity.User), args.String(1), args.String(2), args.Error(3)
}

func (m *MockAuthUsecase) Refresh(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}

func (m *MockAuthUsecase) UpdateProfile(userID uint, name, phone *string, cityID *uint) (*entity.User, error) {
	args := m.Called(userID, name, phone, cityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthUsecase) GetCurrentUser(userID uint) (*entity.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthUsecase) ChangePassword(userID uint, currentPassword, newPassword string) error {
	args := m.Called(userID, currentPassword, newPassword)
	return args.Error(0)
}

func (m *MockAuthUsecase) SetPasswordFromInvite(token, password string) (*entity.User, string, string, error) {
	args := m.Called(token, password)
	if args.Get(0) == nil {
		return nil, "", "", args.Error(3)
	}
	return args.Get(0).(*entity.User), args.String(1), args.String(2), args.Error(3)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthUsecase)
		expectedStatus int
	}{
		{
			name: "successful registration",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
				"name":     "Test User",
			},
			mockSetup: func(m *MockAuthUsecase) {
				user := &entity.User{
					ID:    1,
					Email: "test@example.com",
					Name:  stringPtr("Test User"),
				}
				m.On("Register", "test@example.com", "password123", mock.Anything, mock.Anything, mock.Anything).Return(user, "token123", "refresh123", nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid email",
			requestBody: map[string]interface{}{
				"email":    "invalid-email",
				"password": "password123",
			},
			mockSetup:      func(m *MockAuthUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			mockSetup:      func(m *MockAuthUsecase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockAuthUsecase)
			tt.mockSetup(mockUsecase)

			handler := NewAuthHandler(mockUsecase)
			router := setupRouter()
			router.POST("/register", handler.Register)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthUsecase)
		expectedStatus int
	}{
		{
			name: "successful login",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			mockSetup: func(m *MockAuthUsecase) {
				user := &entity.User{
					ID:    1,
					Email: "test@example.com",
				}
				m.On("Login", "test@example.com", "password123").Return(user, "token123", "refresh123", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			mockSetup: func(m *MockAuthUsecase) {
				m.On("Login", "test@example.com", "wrongpassword").Return(nil, "", "", assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockAuthUsecase)
			tt.mockSetup(mockUsecase)

			handler := NewAuthHandler(mockUsecase)
			router := setupRouter()
			router.POST("/login", handler.Login)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
