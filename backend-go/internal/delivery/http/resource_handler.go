package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
)

type ResourceHandler struct {
	tableRepo       domain.TableRepository
	sessionTypeRepo domain.DiningSessionTypeRepository
}

func NewResourceHandler(t domain.TableRepository, st domain.DiningSessionTypeRepository) *ResourceHandler {
	return &ResourceHandler{
		tableRepo:       t,
		sessionTypeRepo: st,
	}
}

func (h *ResourceHandler) ListTables(c *gin.Context) {
	tables, err := h.tableRepo.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tables)
}

func (h *ResourceHandler) ListSessionTypes(c *gin.Context) {
	types, err := h.sessionTypeRepo.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types)
}

func (h *ResourceHandler) RegisterProtectedRoutes(router gin.IRoutes) {
	router.GET("/tables", h.ListTables)
	router.GET("/session-types", h.ListSessionTypes)
}
