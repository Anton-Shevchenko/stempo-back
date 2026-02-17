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

type MockBusinessUsecase struct {
	mock.Mock
}

func (m *MockBusinessUsecase) Create(business *entity.Business) error {
	args := m.Called(business)
	return args.Error(0)
}

func (m *MockBusinessUsecase) GetByID(id uint) (*entity.Business, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Business), args.Error(1)
}

func (m *MockBusinessUsecase) GetByOwnerID(ownerID uint) (*entity.Business, error) {
	args := m.Called(ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Business), args.Error(1)
}

func (m *MockBusinessUsecase) GetFeatured() ([]entity.Business, error) {
	args := m.Called()
	return args.Get(0).([]entity.Business), args.Error(1)
}

func (m *MockBusinessUsecase) GetAll(page, pageSize int) ([]entity.Business, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]entity.Business), args.Get(1).(int64), args.Error(2)
}

func (m *MockBusinessUsecase) Update(business *entity.Business, userID uint) error {
	args := m.Called(business, userID)
	return args.Error(0)
}

func (m *MockBusinessUsecase) Delete(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func TestBusinessHandler_GetByID(t *testing.T) {
	mockUsecase := new(MockBusinessUsecase)
	business := &entity.Business{
		ID:      1,
		Name:    "Test Business",
		Address: "123 Main St",
	}

	mockUsecase.On("GetByID", uint(1)).Return(business, nil)

	handler := NewBusinessHandler(mockUsecase)
	router := gin.New()
	router.GET("/:id", handler.GetByID)

	req := httptest.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestBusinessHandler_Create(t *testing.T) {
	mockUsecase := new(MockBusinessUsecase)
	business := &entity.Business{
		Name:     "New Business",
		Category: "coffee",
		Address:  "456 Oak St",
		OwnerID:  1,
	}

	mockUsecase.On("Create", mock.AnythingOfType("*entity.Business")).Return(nil)

	handler := NewBusinessHandler(mockUsecase)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
	})
	router.POST("", handler.Create)

	body, _ := json.Marshal(business)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestBusinessHandler_GetAll(t *testing.T) {
	mockUsecase := new(MockBusinessUsecase)
	businesses := []entity.Business{
		{ID: 1, Name: "Business 1"},
		{ID: 2, Name: "Business 2"},
	}
	total := int64(2)

	mockUsecase.On("GetAll", 1, 10).Return(businesses, total, nil)

	handler := NewBusinessHandler(mockUsecase)
	router := gin.New()
	router.GET("", handler.GetAll)

	req := httptest.NewRequest("GET", "/?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestBusinessHandler_Update(t *testing.T) {
	mockUsecase := new(MockBusinessUsecase)
	business := &entity.Business{
		ID:      1,
		Name:    "Updated Business",
		Address: "789 Pine St",
	}

	mockUsecase.On("Update", mock.AnythingOfType("*entity.Business"), uint(1)).Return(nil)

	handler := NewBusinessHandler(mockUsecase)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
	})
	router.PUT("/:id", handler.Update)

	body, _ := json.Marshal(business)
	req := httptest.NewRequest("PUT", "/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestBusinessHandler_Delete(t *testing.T) {
	mockUsecase := new(MockBusinessUsecase)

	mockUsecase.On("Delete", uint(1), uint(1)).Return(nil)

	handler := NewBusinessHandler(mockUsecase)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
	})
	router.DELETE("/:id", handler.Delete)

	req := httptest.NewRequest("DELETE", "/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}
