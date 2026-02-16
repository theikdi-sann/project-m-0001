package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/theikdi-sann/qr-restaurant-api/internal/domain"
	"github.com/theikdi-sann/qr-restaurant-api/internal/usecase"
)

type OrderHandler struct {
	usecase usecase.OrderUsecase
}

func NewOrderHandler(u usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		usecase: u,
	}
}

// DTOs
type createOrderItemRequest struct {
	MenuItemID string `json:"menu_item_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	Notes      string `json:"notes"`
}

type createOrderRequest struct {
	DiningSessionID string                   `json:"dining_session_id" binding:"required"`
	Items           []createOrderItemRequest `json:"items" binding:"required,min=1"`
}

type updateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// Handlers

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := uuid.Parse(req.DiningSessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Session ID"})
		return
	}

	var items []usecase.CreateOrderItemInput
	for _, itemReq := range req.Items {
		menuID, err := uuid.Parse(itemReq.MenuItemID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Menu Item ID"})
			return
		}
		items = append(items, usecase.CreateOrderItemInput{
			MenuItemID: menuID,
			Quantity:   itemReq.Quantity,
			Notes:      itemReq.Notes,
		})
	}

	input := usecase.CreateOrderInput{
		DiningSessionID: sessionID,
		Items:           items,
	}

	order, err := h.usecase.CreateOrder(c.Request.Context(), input)
	if err != nil {
		switch err {
		case domain.ErrSessionNotActive, domain.ErrSessionExpired:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case domain.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Session or Item not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Session ID"})
		return
	}

	orders, err := h.usecase.ListOrders(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID"})
		return
	}

	var req updateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := domain.OrderStatus(req.Status)
	// Basic validation of status string presence is handled by domain logic, 
	// but we could add a quick check here if we wanted.

	order, err := h.usecase.UpdateOrderStatus(c.Request.Context(), orderID, status)
	if err != nil {
		switch err {
		case domain.ErrInvalidStatusTransition:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) RegisterRoutes(router gin.IRoutes) {
	router.POST("/orders", h.CreateOrder)
	router.GET("/orders", h.ListOrders)
	router.PATCH("/orders/:id/status", h.UpdateStatus)
}
