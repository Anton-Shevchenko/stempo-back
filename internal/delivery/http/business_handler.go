package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/infrastructure/middleware"
	"github.com/stempo/backend/internal/usecase"
)

type BusinessHandler struct {
	businessUsecase usecase.BusinessUsecase
}

func NewBusinessHandler(businessUsecase usecase.BusinessUsecase) *BusinessHandler {
	return &BusinessHandler{businessUsecase: businessUsecase}
}

func (h *BusinessHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var business entity.Business
	if err := c.ShouldBindJSON(&business); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Clear ID and other fields that should not be set on create
	business.ID = 0
	business.OwnerID = userID
	business.Status = entity.BusinessStatusPending
	business.Rating = 0
	business.IsOpen = true
	business.HasLoyaltyProgram = business.HasLoyaltyProgram || false
	business.Featured = false
	
	if err := h.businessUsecase.Create(&business); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, business)
}

func (h *BusinessHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	business, err := h.businessUsecase.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "business not found"})
		return
	}

	c.JSON(http.StatusOK, business)
}

func (h *BusinessHandler) GetFeatured(c *gin.Context) {
	businesses, err := h.businessUsecase.GetFeatured()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, businesses)
}

func (h *BusinessHandler) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
		return
	}

	businesses, err := h.businessUsecase.GetByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, businesses)
}

func (h *BusinessHandler) GetNewest(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	businesses, err := h.businessUsecase.GetNewest(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, businesses)
}

func (h *BusinessHandler) GetMyBusiness(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	business, err := h.businessUsecase.GetByOwnerID(userID)
	if err != nil {
		// If business not found, return null instead of error
		c.JSON(http.StatusOK, nil)
		return
	}

	c.JSON(http.StatusOK, business)
}

func (h *BusinessHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	businesses, total, err := h.businessUsecase.GetAll(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":    businesses,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *BusinessHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var business entity.Business
	if err := c.ShouldBindJSON(&business); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	business.ID = uint(id)
	// Don't allow changing OwnerID, Status, Rating, CreatedAt, UpdatedAt via update
	// These fields are managed by the system and will be preserved in usecase
	
	if err := h.businessUsecase.Update(&business, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, business)
}

func (h *BusinessHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.businessUsecase.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "business deleted successfully"})
}
