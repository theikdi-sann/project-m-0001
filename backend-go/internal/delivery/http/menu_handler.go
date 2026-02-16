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

	items, err := h.menuRepo.ListByCategory(c.Request.Context(), catID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *MenuHandler) RegisterPublicRoutes(router gin.IRoutes) {
	router.GET("/menu/items", h.GetItemsByCategory)
	router.GET("/menu/categories", h.GetCategories)
}
