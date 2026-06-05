package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/infrastructure/middleware"
	"github.com/stempo/backend/internal/usecase"
)

type BonusProgramHandler struct {
	programUsecase  usecase.BonusProgramUsecase
	businessUsecase usecase.BusinessUsecase
}

func NewBonusProgramHandler(programUsecase usecase.BonusProgramUsecase, businessUsecase usecase.BusinessUsecase) *BonusProgramHandler {
	return &BonusProgramHandler{
		programUsecase:  programUsecase,
		businessUsecase: businessUsecase,
	}
}

// Helper function to get and validate user ID from context
func (h *BonusProgramHandler) getUserID(c *gin.Context) (uint, bool) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return 0, false
	}
	return userID, true
}

// Helper function to parse uint ID from URL parameter
func parseUintParam(c *gin.Context, paramName string) (uint, bool) {
	idStr := c.Param(paramName)
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": paramName + " is required"})
		return 0, false
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + paramName})
		return 0, false
	}

	return uint(id), true
}

// Helper function to verify business ownership
func (h *BonusProgramHandler) verifyBusinessOwnership(c *gin.Context, businessID uint, userID uint) bool {
	business, err := h.businessUsecase.GetByID(businessID)
	if err != nil || business == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "business not found"})
		return false
	}

	if business.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: business does not belong to you"})
		return false
	}

	return true
}

func (h *BonusProgramHandler) Create(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var program entity.BonusProgram
	if err := c.ShouldBindJSON(&program); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.programUsecase.Create(&program, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, program)
}

func (h *BonusProgramHandler) GetByID(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	program, err := h.programUsecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "program not found"})
		return
	}

	c.JSON(http.StatusOK, program)
}

// GetPublicByBusinessID returns approved/active programs for a business (no auth required).
// Route: GET /api/bonus-programs/public/by-business/:id
func (h *BonusProgramHandler) GetPublicByBusinessID(c *gin.Context) {
	businessID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	programs, err := h.programUsecase.GetActiveByBusinessID(businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if programs == nil {
		programs = []entity.BonusProgram{}
	}

	c.JSON(http.StatusOK, programs)
}

// GetByBusinessID returns all programs for a business
// Requires authentication and verifies that the business belongs to the requesting user
// Route: GET /api/bonus-programs/by-business/:id
func (h *BonusProgramHandler) GetByBusinessID(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	businessID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	if !h.verifyBusinessOwnership(c, businessID, userID) {
		return
	}

	programs, err := h.programUsecase.GetByBusinessID(businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, programs)
}

func (h *BonusProgramHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	programs, total, err := h.programUsecase.GetAll(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":    programs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *BonusProgramHandler) Update(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var program entity.BonusProgram
	if err := c.ShouldBindJSON(&program); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	program.ID = id

	if err := h.programUsecase.Update(&program, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, program)
}

func (h *BonusProgramHandler) Delete(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.programUsecase.Delete(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "program deleted successfully"})
}
