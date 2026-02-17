package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/infrastructure/middleware"
	"github.com/stempo/backend/internal/usecase"
)

type LoyaltyCardHandler struct {
	cardUsecase usecase.LoyaltyCardUsecase
}

func NewLoyaltyCardHandler(cardUsecase usecase.LoyaltyCardUsecase) *LoyaltyCardHandler {
	return &LoyaltyCardHandler{cardUsecase: cardUsecase}
}

func (h *LoyaltyCardHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		BusinessID uint `json:"businessId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card := &entity.LoyaltyCard{
		UserID:         userID,
		BusinessID:     req.BusinessID,
		Stamps:         0,
		StampsRequired:  10,
	}

	if err := h.cardUsecase.Create(card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, card)
}

func (h *LoyaltyCardHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	card, err := h.cardUsecase.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}

	c.JSON(http.StatusOK, card)
}

func (h *LoyaltyCardHandler) GetByUserID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cards, err := h.cardUsecase.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func (h *LoyaltyCardHandler) AddStamp(c *gin.Context) {
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

	if err := h.cardUsecase.AddStamp(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, _ := h.cardUsecase.GetByID(uint(id))
	c.JSON(http.StatusOK, card)
}

func (h *LoyaltyCardHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var card entity.LoyaltyCard
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card.ID = uint(id)
	if err := h.cardUsecase.Update(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}
