package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct{}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

// @Summary Recebe notificações (endpoint de teste)
// @Description Endpoint para testar o recebimento de notificações genéricas
// @Tags Receptor
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param notification body object true "JSON genérico de notificação"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} Response
// @Router /api/webhook/test/notifications [post]
func (h *WebhookHandler) ReceiveNotification(c *gin.Context) {
	var notification map[string]interface{}
	if err := c.ShouldBindJSON(&notification); err != nil {
		log.Printf("Erro ao decodificar JSON: %v", err)
		c.JSON(http.StatusBadRequest, Response{
			Error: "dados inválidos",
		})
		return
	}

	log.Printf("Notificação recebida: %+v", notification)

	c.JSON(http.StatusOK, Response{
		Message: "Notificação recebida com sucesso",
		Data:    notification,
	})
}

func (h *WebhookHandler) SetupRoutes(r *gin.RouterGroup) {
	webhook := r.Group("/webhook/test")
	{
		webhook.POST("/notifications", h.ReceiveNotification)
	}
}
