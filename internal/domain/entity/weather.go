package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type WeatherForecast struct {
	Date     time.Time
	MinTemp  float64
	MaxTemp  float64
	Forecast string
	UV       float64
	Wave     *WaveInfo
}

type WaveInfo struct {
	UpdateTime string     `xml:"atualizacao"`
	Morning    WavePeriod `xml:"manha"`
	Afternoon  WavePeriod `xml:"tarde"`
	Night      WavePeriod `xml:"noite"`
}

type WavePeriod struct {
	Date      string  `xml:"dia"`
	Agitation string  `xml:"agitacao"`
	Height    float64 `xml:"altura"`
	Direction string  `xml:"direcao"`
	WindSpeed float64 `xml:"vento"`
	WindDir   string  `xml:"vento_dir"`
}

type WeatherForecastCollection struct {
	LocationID uuid.UUID
	Nome       string
	UF         string
	Forecasts  []WeatherForecast
	UpdatedAt  time.Time
}

func NewWeatherForecastCollection(locationID uuid.UUID, nome, uf string, forecasts []WeatherForecast) *WeatherForecastCollection {
	return &WeatherForecastCollection{
		LocationID: locationID,
		Nome:       nome,
		UF:         uf,
		Forecasts:  forecasts,
		UpdatedAt:  time.Now(),
	}
}

func (w *WeatherForecastCollection) GetNext4Days() []WeatherForecast {
	if len(w.Forecasts) <= 4 {
		return w.Forecasts
	}
	return w.Forecasts[:4]
}

func (w *WeatherForecast) HasWaveForecast() bool {
	return w.Wave != nil
}

func (w *WeatherForecast) FormatTemperature() string {
	return fmt.Sprintf("%.1f°C / %.1f°C", w.MinTemp, w.MaxTemp)
}

func (w *WeatherForecast) AsNotificationText() string {
	text := fmt.Sprintf("%s: %s - %s",
		w.Date.Format("02/01"),
		w.FormatTemperature(),
		w.Forecast,
	)

	if w.HasWaveForecast() {
		text += fmt.Sprintf(" | Ondas - Manhã: %.1fm %s, Vento: %.1f km/h %s",
			w.Wave.Morning.Height,
			w.Wave.Morning.Direction,
			w.Wave.Morning.WindSpeed,
			w.Wave.Morning.WindDir,
		)

		text += fmt.Sprintf(" | Tarde: %.1fm %s, Vento: %.1f km/h %s",
			w.Wave.Afternoon.Height,
			w.Wave.Afternoon.Direction,
			w.Wave.Afternoon.WindSpeed,
			w.Wave.Afternoon.WindDir,
		)

		text += fmt.Sprintf(" | Noite: %.1fm %s, Vento: %.1f km/h %s",
			w.Wave.Night.Height,
			w.Wave.Night.Direction,
			w.Wave.Night.WindSpeed,
			w.Wave.Night.WindDir,
		)
	}

	return text
}
