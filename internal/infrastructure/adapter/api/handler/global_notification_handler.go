package handler

import (
	"net/http"
	"time"
	"weather-notification/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type GlobalNotificationHandler struct {
	globalNotificationService *service.GlobalNotificationService
}

func NewGlobalNotificationHandler(globalNotificationService *service.GlobalNotificationService) *GlobalNotificationHandler {
	return &GlobalNotificationHandler{
		globalNotificationService: globalNotificationService,
	}
}

// @Summary Cria uma notificação global
// @Description Cria uma notificação que será enviada para todos os usuários ativos no horário especificado
// @Tags Notificações Globais
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateGlobalNotificationRequest true "Dados da notificação global"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/notifications/global [post]
func (h *GlobalNotificationHandler) Create(c *gin.Context) {
	var req CreateGlobalNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "dados inválidos: " + err.Error(),
		})
		return
	}

	timeOfDay, err := time.Parse("15:04", req.TimeOfDay)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "formato de hora inválido",
		})
		return
	}

	err = h.globalNotificationService.Create(c.Request.Context(), timeOfDay, req.Frequency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Message: "Notificação global criada com sucesso",
	})
}

// @Summary Lista notificações globais ativas
// @Description Retorna todas as notificações globais ativas
// @Tags Notificações Globais
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/notifications/global [get]
func (h *GlobalNotificationHandler) List(c *gin.Context) {
	notifications, err := h.globalNotificationService.ListActive(c.Request.Context())
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

func (h *GlobalNotificationHandler) SetupRoutes(r *gin.RouterGroup) {
	notifications := r.Group("/notifications/global")
	{
		notifications.POST("", h.Create)
		notifications.GET("", h.List)
	}
}
