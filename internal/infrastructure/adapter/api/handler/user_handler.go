package handler

import (
	"net/http"
	"strings"
	"unicode"
	"weather-notification/internal/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/text/runes"
	"golang.org/x/text/unicode/norm"
)

type UserHandler struct {
	userService    *service.UserService
	weatherService *service.WeatherService
}

func NewUserHandler(userService *service.UserService, weatherService *service.WeatherService) *UserHandler {
	return &UserHandler{
		userService:    userService,
		weatherService: weatherService,
	}
}

// @Summary Cria um novo usuário
// @Description Cria um usuário e vincula a uma localização
// @Tags Usuários
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Dados do usuário"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "dados inválidos: " + err.Error(),
		})
		return
	}

	req.City = strings.Map(func(r rune) rune {
		if runes.In(unicode.Mn).Contains(r) {
			return -1
		}
		return r
	}, norm.NFD.String(req.City))

	locations, err := h.weatherService.SearchLocation(c.Request.Context(), req.City)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: "erro ao buscar localização: " + err.Error(),
		})
		return
	}

	err = h.userService.Create(c.Request.Context(), req.Name, req.Email, locations[0].ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Message: "Usuário criado com sucesso",
	})
}

// @Summary Atualiza um usuário
// @Description Atualiza os dados de um usuário
// @Tags Usuários
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário" Format(uuid)
// @Param request body UpdateUserRequest true "Dados para atualização"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "ID inválido",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "dados inválidos: " + err.Error(),
		})
		return
	}

	var locationID *uuid.UUID
	if req.City != "" {
		location, err := h.weatherService.SearchLocation(c.Request.Context(), req.City)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Error: "erro ao buscar localização: " + err.Error(),
			})
			return
		}
		locationID = &location[0].ID
	}

	err = h.userService.Update(c.Request.Context(), userID, req.Name, *locationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Message: "Usuário atualizado com sucesso",
	})
}

// @Summary Lista usuários
// @Description Retorna uma lista de todos os usuários cadastrados
// @Tags Usuários
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /api/users [get]
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: users,
	})
}

// @Summary Ativa ou desativa o opt-out do usuário
// @Description Permite ao usuário ativar ou desativar o opt-out de notificações
// @Tags Usuários
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path string true "ID do usuário" Format(uuid)
// @Param request body ToggleOptOutRequest true "Novo status do opt-out"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/users/{user_id}/optout [patch]
func (h *UserHandler) ToggleOptOut(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "user_id inválido",
		})
		return
	}

	var req ToggleOptOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "dados inválidos: " + err.Error(),
		})
		return
	}

	// Atualiza o status de opt-out do usuário
	err = h.userService.ToggleOptOut(c.Request.Context(), userID, req.OptOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Message: "Status de opt-out atualizado com sucesso",
	})
}

func (h *UserHandler) SetupRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("", h.Create)
		users.PUT("/:id", h.Update)
		users.GET("", h.List)
		users.PATCH("/:user_id/optout", h.ToggleOptOut)
	}
}
