package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stempo/backend/internal/infrastructure/middleware"
	"github.com/stempo/backend/internal/usecase"
)

type EmployeeHandler struct {
	employeeUsecase usecase.EmployeeUsecase
}

func NewEmployeeHandler(employeeUsecase usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{employeeUsecase: employeeUsecase}
}

type AddEmployeeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *EmployeeHandler) AddEmployee(c *gin.Context) {
	ownerID := middleware.GetUserID(c)
	if ownerID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	businessID, err := strconv.ParseUint(c.Param("businessId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business id"})
		return
	}

	var req AddEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := h.employeeUsecase.AddEmployee(uint(businessID), ownerID, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

func (h *EmployeeHandler) RemoveEmployee(c *gin.Context) {
	ownerID := middleware.GetUserID(c)
	if ownerID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	businessID, err := strconv.ParseUint(c.Param("businessId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business id"})
		return
	}

	employeeID, err := strconv.ParseUint(c.Param("employeeId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	if err := h.employeeUsecase.RemoveEmployee(uint(businessID), uint(employeeID), ownerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee removed successfully"})
}

func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	ownerID := middleware.GetUserID(c)
	if ownerID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	businessID, err := strconv.ParseUint(c.Param("businessId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business id"})
		return
	}

	employees, err := h.employeeUsecase.GetEmployeesByBusiness(uint(businessID), ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) GetMyBusinesses(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	businesses, err := h.employeeUsecase.GetBusinessesByEmployee(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, businesses)
}
