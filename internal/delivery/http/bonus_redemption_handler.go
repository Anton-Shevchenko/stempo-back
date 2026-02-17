package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stempo/backend/internal/infrastructure/middleware"
	"github.com/stempo/backend/internal/usecase"
)

type BonusRedemptionHandler struct {
	redemptionUsecase usecase.BonusRedemptionUsecase
}

func NewBonusRedemptionHandler(redemptionUsecase usecase.BonusRedemptionUsecase) *BonusRedemptionHandler {
	return &BonusRedemptionHandler{redemptionUsecase: redemptionUsecase}
}

func (h *BonusRedemptionHandler) GenerateCode(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cardID, err := strconv.ParseUint(c.Param("cardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	redemption, err := h.redemptionUsecase.GenerateRedemptionCode(uint(cardID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, redemption)
}

func (h *BonusRedemptionHandler) Redeem(c *gin.Context) {
	scannedByUserID := middleware.GetUserID(c)
	if scannedByUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redemption, err := h.redemptionUsecase.RedeemBonus(req.Code, scannedByUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, redemption)
}

func (h *BonusRedemptionHandler) GetByCard(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cardID, err := strconv.ParseUint(c.Param("cardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	redemptions, err := h.redemptionUsecase.GetRedemptionsByCard(uint(cardID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, redemptions)
}
