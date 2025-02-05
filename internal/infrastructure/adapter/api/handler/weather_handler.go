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

type WeatherHandler struct {
	weatherService *service.WeatherService
}

func NewWeatherHandler(weatherService *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

// @Summary Busca cidade por nome
// @Description Busca uma cidade no CPTEC por nome
// @Tags Localizações
// @Security BearerAuth
// @Produce json
// @Param city query string true "Nome da cidade"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/weather/search [get]
func (h *WeatherHandler) SearchLocation(c *gin.Context) {
	cityName := c.Query("city")
	if cityName == "" {
		c.JSON(http.StatusBadRequest, Response{
			Error: "nome da cidade é obrigatório",
		})
		return
	}

	cityName = strings.Map(func(r rune) rune {
		if runes.In(unicode.Mn).Contains(r) {
			return -1
		}
		return r
	}, norm.NFD.String(cityName))

	locations, err := h.weatherService.SearchLocation(c.Request.Context(), cityName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: locations,
	})
}

// @Summary Busca previsão do tempo
// @Description Retorna a previsão do tempo para uma localidade
// @Tags Clima
// @Security BearerAuth
// @Produce json
// @Param location_id query string true "ID da localidade" Format(uuid)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/weather/forecast [get]
func (h *WeatherHandler) GetForecast(c *gin.Context) {
	locationID, err := uuid.Parse(c.Query("location_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: "location_id inválido",
		})
		return
	}

	forecast, err := h.weatherService.GetForecast(c.Request.Context(), locationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: forecast,
	})
}

func (h *WeatherHandler) SetupRoutes(r *gin.RouterGroup) {
	weather := r.Group("/weather")
	{
		weather.GET("/search", h.SearchLocation)
		weather.GET("/forecast", h.GetForecast)
	}
}
