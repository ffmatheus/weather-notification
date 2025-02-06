package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type WeatherForecast struct {
	Date     time.Time `json:"date"`
	MinTemp  float64   `json:"min_temp"`
	MaxTemp  float64   `json:"max_temp"`
	Forecast string    `json:"forecast"`
	UV       float64   `json:"uv"`
	Wave     *WaveInfo `json:"wave,omitempty"`
}

type WaveInfo struct {
	UpdateTime string     `xml:"atualizacao" json:"update_time"`
	Morning    WavePeriod `xml:"manha" json:"morning"`
	Afternoon  WavePeriod `xml:"tarde" json:"afternoon"`
	Night      WavePeriod `xml:"noite" json:"night"`
}

type WavePeriod struct {
	Date      string  `xml:"dia" json:"date"`
	Agitation string  `xml:"agitacao" json:"agitation"`
	Height    float64 `xml:"altura" json:"height"`
	Direction string  `xml:"direcao" json:"direction"`
	WindSpeed float64 `xml:"vento" json:"wind_speed"`
	WindDir   string  `xml:"vento_dir" json:"wind_dir"`
}

type WeatherForecastCollection struct {
	Nome      string            `json:"nome"`
	UF        string            `json:"uf"`
	Forecasts []WeatherForecast `json:"forecasts"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func NewWeatherForecastCollection(locationID uuid.UUID, nome, uf string, forecasts []WeatherForecast) *WeatherForecastCollection {
	return &WeatherForecastCollection{
		Nome:      nome,
		UF:        uf,
		Forecasts: forecasts,
		UpdatedAt: time.Now(),
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
