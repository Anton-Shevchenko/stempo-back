package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stempo/backend/internal/usecase"
)

type CityHandler struct {
	cityUsecase usecase.CityUsecase
}

func NewCityHandler(cityUsecase usecase.CityUsecase) *CityHandler {
	return &CityHandler{cityUsecase: cityUsecase}
}

func (h *CityHandler) GetAll(c *gin.Context) {
	cities, err := h.cityUsecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cities)
}

func (h *CityHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	city, err := h.cityUsecase.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "city not found"})
		return
	}

	c.JSON(http.StatusOK, city)
}

type CreateCityRequest struct {
	Name      string   `json:"name" binding:"required"`
	Slug      string   `json:"slug,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

func (h *CityHandler) Create(c *gin.Context) {
	var req CreateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	city, err := h.cityUsecase.Create(req.Name, req.Slug, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, city)
}

type UpdateCityRequest struct {
	Name      string   `json:"name,omitempty"`
	Slug      string   `json:"slug,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

func (h *CityHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	city, err := h.cityUsecase.Update(uint(id), req.Name, req.Slug, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, city)
}

func (h *CityHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.cityUsecase.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "city deleted successfully"})
}
