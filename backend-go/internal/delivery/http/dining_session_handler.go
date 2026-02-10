package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/theikdi-sann/qr-restaurant-api/internal/usecase"
)

type DiningSessionHandler struct {
	usecase usecase.DiningSessionUsecase
}

func NewDiningSessionHandler(u usecase.DiningSessionUsecase) *DiningSessionHandler {
	return &DiningSessionHandler{
		usecase: u,
	}
}

type createSessionRequest struct {
	TableID       string `json:"table_id" binding:"required"`
	SessionTypeID string `json:"session_type_id" binding:"required"`
	CreatedBy     string `json:"created_by" binding:"required"` // In real app, from Auth Context
	GuestCount    int    `json:"guest_count" binding:"required,min=1"`
}

func (h *DiningSessionHandler) CreateSession(c *gin.Context) {
	var req createSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tableID, err := uuid.Parse(req.TableID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Table ID"})
		return
	}

	sessionTypeID, err := uuid.Parse(req.SessionTypeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Session Type ID"})
		return
	}

	createdBy, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	input := usecase.CreateSessionInput{
		TableID:       tableID,
		SessionTypeID: sessionTypeID,
		CreatedBy:     createdBy,
		GuestCount:    req.GuestCount,
	}

	session, err := h.usecase.CreateSession(c.Request.Context(), input)
	if err != nil {
		// In a real app, check for specific domain errors (e.g. ErrTableOccupied) and return 409 Conflict
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *DiningSessionHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/sessions", h.CreateSession)
}
