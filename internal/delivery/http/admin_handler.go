package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/usecase"
)

type AdminHandler struct {
	businessUsecase usecase.BusinessUsecase
	programUsecase  usecase.BonusProgramUsecase
}

func NewAdminHandler(businessUsecase usecase.BusinessUsecase, programUsecase usecase.BonusProgramUsecase) *AdminHandler {
	return &AdminHandler{
		businessUsecase: businessUsecase,
		programUsecase:  programUsecase,
	}
}

func (h *AdminHandler) GetPendingBusinesses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	businesses, total, err := h.businessUsecase.GetByStatus(entity.BusinessStatusPending, page, pageSize)
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

func (h *AdminHandler) ApproveBusiness(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.businessUsecase.Approve(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "business approved successfully"})
}

func (h *AdminHandler) RejectBusiness(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.businessUsecase.Reject(uint(id), req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "business rejected successfully"})
}

func (h *AdminHandler) GetPendingPrograms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	programs, total, err := h.programUsecase.GetByStatus(entity.BonusProgramStatusPending, page, pageSize)
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

func (h *AdminHandler) ApproveProgram(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.programUsecase.Approve(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bonus program approved successfully"})
}

func (h *AdminHandler) RejectProgram(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.programUsecase.Reject(uint(id), req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bonus program rejected successfully"})
}

func (h *AdminHandler) UpdateBusiness(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Check if business exists
	_, err = h.businessUsecase.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "business not found"})
		return
	}

	// Read raw body and parse JSON into map to check which fields were actually provided
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	var updateData map[string]interface{}
	if err := json.Unmarshal(body, &updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build updates map with only provided fields
	updates := make(map[string]interface{})
	
	if name, ok := updateData["name"].(string); ok {
		updates["name"] = name
	}
	if category, ok := updateData["category"].(string); ok {
		updates["category"] = category
	}
	if address, ok := updateData["address"].(string); ok {
		updates["address"] = address
	}
	if description, ok := updateData["description"]; ok {
		if description == nil {
			updates["description"] = nil
		} else if descStr, ok := description.(string); ok {
			updates["description"] = descStr
		}
	}
	if imageURL, ok := updateData["imageUrl"]; ok {
		if imageURL == nil {
			updates["image_url"] = nil
		} else if urlStr, ok := imageURL.(string); ok {
			updates["image_url"] = urlStr
		}
	}
	if icon, ok := updateData["icon"].(string); ok {
		updates["icon"] = icon
	}
	if iconColor, ok := updateData["iconColor"]; ok {
		if iconColor == nil {
			updates["icon_color"] = nil
		} else if colorStr, ok := iconColor.(string); ok {
			updates["icon_color"] = colorStr
		}
	}
	if featured, ok := updateData["featured"].(bool); ok {
		updates["featured"] = featured
	}
	if status, ok := updateData["status"].(string); ok {
		updates["status"] = entity.BusinessStatus(status)
	}
	if rating, ok := updateData["rating"]; ok {
		if ratingFloat, ok := rating.(float64); ok {
			updates["rating"] = ratingFloat
		}
	}
	if isOpen, ok := updateData["isOpen"].(bool); ok {
		updates["is_open"] = isOpen
	}
	if hasLoyaltyProgram, ok := updateData["hasLoyaltyProgram"].(bool); ok {
		updates["has_loyalty_program"] = hasLoyaltyProgram
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	// Use UpdateFields to update only provided fields
	if err := h.businessUsecase.UpdateFields(uint(id), updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return updated business
	updated, err := h.businessUsecase.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve updated business"})
		return
	}

	c.JSON(http.StatusOK, updated)
}
