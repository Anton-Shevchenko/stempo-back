package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/infrastructure/middleware"
	"github.com/stempo/backend/internal/usecase"
)

type QRCodeHandler struct {
	qrCodeUsecase usecase.QRCodeUsecase
}

func NewQRCodeHandler(qrCodeUsecase usecase.QRCodeUsecase) *QRCodeHandler {
	return &QRCodeHandler{
		qrCodeUsecase: qrCodeUsecase,
	}
}

type GenerateQRRequest struct {
	ProgramID       uint                `json:"programId" binding:"required"`
	Type            entity.QRCodeType   `json:"type" binding:"required"`
	ExpirationHours *int                `json:"expirationHours,omitempty"`
}

func (h *QRCodeHandler) Generate(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req GenerateQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate QR code type
	if req.Type != entity.QRCodeTypeTemporary && req.Type != entity.QRCodeTypePermanent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid QR code type. Must be 'temporary' or 'permanent'"})
		return
	}

	qrCode, err := h.qrCodeUsecase.Generate(req.ProgramID, req.Type, req.ExpirationHours, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, qrCode)
}

func (h *QRCodeHandler) Validate(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QR code is required"})
		return
	}

	qrCode, err := h.qrCodeUsecase.Validate(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, qrCode)
}

func (h *QRCodeHandler) GetByProgramID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	programID, ok := parseUintParam(c, "programId")
	if !ok {
		return
	}

	qrCodes, err := h.qrCodeUsecase.GetByProgramID(programID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, qrCodes)
}

func (h *QRCodeHandler) GetByBusinessID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	businessID, ok := parseUintParam(c, "businessId")
	if !ok {
		return
	}

	qrCodes, err := h.qrCodeUsecase.GetByBusinessID(businessID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, qrCodes)
}

func (h *QRCodeHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid QR code ID"})
		return
	}

	if err := h.qrCodeUsecase.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "QR code deleted successfully"})
}
