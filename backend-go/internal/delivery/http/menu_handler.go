package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
)

type MenuHandler struct {
	menuRepo domain.MenuItemRepository
}

func NewMenuHandler(r domain.MenuItemRepository) *MenuHandler {
	return &MenuHandler{
		menuRepo: r,
	}
}

func (h *MenuHandler) GetCategories(c *gin.Context) {
	categories, err := h.menuRepo.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *MenuHandler) GetItemsByCategory(c *gin.Context) {
	catIDStr := c.Query("category_id")
	if catIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id required"})
		return
	}

	catID, err := uuid.Parse(catIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Category ID"})
		return
	}

	// If staff/auth present, maybe list all?
	// For simplicity, I'll add `all=true` query param or check role.
	// But current requirement is "Item Availability Toggles" -> Staff needs a list.
	// I'll create a separate endpoint `GET /staff/menu/items` or reusing with query param.
	// Let's use `include_unavailable=true` query param.
	
	includeUnavailable := c.Query("include_unavailable") == "true"

	var items []*domain.MenuItem
	var errRepo error

	if includeUnavailable {
		items, errRepo = h.menuRepo.ListAllByCategory(c.Request.Context(), catID)
	} else {
		items, errRepo = h.menuRepo.ListByCategory(c.Request.Context(), catID)
	}

	if errRepo != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errRepo.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

type updateAvailabilityRequest struct {
	IsAvailable bool `json:"is_available"`
}

func (h *MenuHandler) UpdateAvailability(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
		return
	}

	var req updateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.menuRepo.UpdateAvailability(c.Request.Context(), id, req.IsAvailable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *MenuHandler) RegisterPublicRoutes(router gin.IRoutes) {
	router.GET("/menu/items", h.GetItemsByCategory)
	router.GET("/menu/categories", h.GetCategories)
}

func (h *MenuHandler) RegisterProtectedRoutes(router gin.IRoutes) {
	router.PATCH("/menu/items/:id/availability", h.UpdateAvailability)
}
