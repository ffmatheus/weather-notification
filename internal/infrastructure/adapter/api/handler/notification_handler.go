package handler

import (
	"net/http"
	"time"
	"weather-notification/internal/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

type CreateNotificationRequest struct {
	UserID      string    `json:"user_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	LocationID  string    `json:"location_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	ScheduleFor time.Time `json:"schedule_for" binding:"required" example:"2025-02-03T21:35:00-03:00"`
}

// @Summary Agenda uma nova notificação
// @Description Cria um agendamento de notificação de previsão do tempo
// @Tags Notificações
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateNotificationRequest true "Dados do agendamento"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/notifications [post]
func (h *NotificationHandler) Create(c *gin.Context) {
	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "dados inválidos: " + err.Error(),
		})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "user_id inválido",
		})
		return
	}

	locationID, err := uuid.Parse(req.LocationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "location_id inválido",
		})
		return
	}

	err = h.notificationService.Schedule(c.Request.Context(), userID, locationID, req.ScheduleFor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Message: "Notificação agendada com sucesso",
	})
}

// @Summary Lista notificações do usuário
// @Description Retorna todas as notificações de um usuário específico
// @Tags Notificações
// @Security BearerAuth
// @Produce json
// @Param user_id query string true "ID do usuário" Format(uuid)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "user_id inválido ou não fornecido",
		})
		return
	}

	notifications, err := h.notificationService.GetUserNotifications(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: notifications,
	})
}

func (h *NotificationHandler) SetupRoutes(r *gin.RouterGroup) {
	notifications := r.Group("/notifications")
	{
		notifications.POST("", h.Create)
		notifications.GET("", h.List)
	}
}
