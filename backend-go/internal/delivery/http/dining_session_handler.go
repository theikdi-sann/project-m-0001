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

	// Get UserID from Context (set by AuthMiddleware)
	userIDVal, exists := c.Get("userID") // Hardcoded key for now, or import from middleware if cyclic dep not an issue
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	input := usecase.CreateSessionInput{
		TableID:       tableID,
		SessionTypeID: sessionTypeID,
		CreatedBy:     userID,
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

func (h *DiningSessionHandler) GetSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Session ID"})
		return
	}

	session, err := h.usecase.GetSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *DiningSessionHandler) RegisterRoutes(router gin.IRoutes) {
	router.POST("/sessions", h.CreateSession)
}

func (h *DiningSessionHandler) RegisterPublicRoutes(router gin.IRoutes) {
	router.GET("/sessions/:id", h.GetSession)
}
