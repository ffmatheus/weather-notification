package cptec

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"weather-notification/internal/domain/entity"

	"github.com/google/uuid"
	"golang.org/x/net/html/charset"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type cityResponse struct {
	XMLName xml.Name `xml:"cidades"`
	Cities  []struct {
		ID    int    `xml:"id"`
		Name  string `xml:"nome"`
		State string `xml:"uf"`
	} `xml:"cidade"`
}

type forecastResponse struct {
	XMLName   xml.Name `xml:"cidade"`
	Name      string   `xml:"nome"`
	State     string   `xml:"uf"`
	Forecasts []struct {
		Date     string `xml:"dia"`
		MinTemp  string `xml:"minima"`
		MaxTemp  string `xml:"maxima"`
		Forecast string `xml:"tempo"`
		IUV      string `xml:"iuv"`
	} `xml:"previsao"`
}

type waveResponse struct {
	XMLName    xml.Name   `xml:"cidade"`
	UpdateTime string     `xml:"atualizacao"`
	Morning    wavePeriod `xml:"manha"`
	Afternoon  wavePeriod `xml:"tarde"`
	Night      wavePeriod `xml:"noite"`
}

type wavePeriod struct {
	Date      string `xml:"dia"`
	Agitation string `xml:"agitacao"`
	Height    string `xml:"altura"`
	Direction string `xml:"direcao"`
	WindSpeed string `xml:"vento"`
	WindDir   string `xml:"vento_dir"`
}

func NewClient() *Client {
	return &Client{
		baseURL: os.Getenv("CPTEC_BASE_URL"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SearchCities(ctx context.Context, cityName string) ([]entity.Location, error) {
	endpoint := fmt.Sprintf("%s/listaCidades?city=%s", c.baseURL, url.QueryEscape(cityName))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code inválido: %d", resp.StatusCode)
	}

	var result cityResponse
	if err := decodeISO88591XML(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar XML: %w", err)
	}

	locations := make([]entity.Location, 0, len(result.Cities))
	for _, city := range result.Cities {
		location, err := entity.NewLocation(
			city.ID,
			city.Name,
			city.State,
		)
		if err != nil {
			continue
		}
		locations = append(locations, *location)
	}

	return locations, nil
}

func (c *Client) GetWeatherForecast(ctx context.Context, cptecCode int) (*entity.WeatherForecastCollection, error) {
	endpoint := fmt.Sprintf("%s/cidade/%d/previsao.xml", c.baseURL, cptecCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code inválido: %d", resp.StatusCode)
	}

	var result forecastResponse
	if err := decodeISO88591XML(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar XML: %w", err)
	}

	result.Name = url.QueryEscape(result.Name)
	result.State = url.QueryEscape(result.State)

	forecasts := make([]entity.WeatherForecast, 0, len(result.Forecasts))
	for _, f := range result.Forecasts {
		date, err := time.Parse("2006-01-02", f.Date)
		if err != nil {
			continue
		}

		iuv, err := strconv.ParseFloat(f.IUV, 64)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter: %w", err)
		}

		minTemp := parseTemperature(f.MinTemp)
		maxTemp := parseTemperature(f.MaxTemp)

		forecast := entity.WeatherForecast{
			Date:     date,
			MinTemp:  minTemp,
			MaxTemp:  maxTemp,
			Forecast: f.Forecast,
			UV:       iuv,
		}
		forecasts = append(forecasts, forecast)
	}

	return entity.NewWeatherForecastCollection(uuid.New(), result.Name, result.State, forecasts), nil
}

func (c *Client) GetWaveForecast(ctx context.Context, cptecCode int, date time.Time) (*entity.WaveInfo, error) {
	endpoint := fmt.Sprintf("%s/cidade/%d/dia/%s/ondas.xml",
		c.baseURL,
		cptecCode,
		date.Format("20060102"),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code inválido: %d", resp.StatusCode)
	}

	var result waveResponse
	if err := decodeISO88591XML(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar XML: %w", err)
	}

	if isEmptyWavePeriod(result.Morning) && isEmptyWavePeriod(result.Afternoon) && isEmptyWavePeriod(result.Night) {
		return nil, nil
	}

	return &entity.WaveInfo{
		UpdateTime: result.UpdateTime,
		Morning:    buildWavePeriod(result.Morning),
		Afternoon:  buildWavePeriod(result.Afternoon),
		Night:      buildWavePeriod(result.Night),
	}, nil
}

func isEmptyWavePeriod(p wavePeriod) bool {
	return p.Height == "0" || p.Height == "" || p.Direction == "" || p.Direction == "undefined"
}

func buildWavePeriod(p wavePeriod) entity.WavePeriod {
	return entity.WavePeriod{
		Date:      p.Date,
		Agitation: p.Agitation,
		Height:    parseTemperature(p.Height),
		Direction: p.Direction,
		WindSpeed: parseTemperature(p.WindSpeed),
		WindDir:   p.WindDir,
	}
}

func decodeISO88591XML(body io.Reader, v interface{}) error {
	decoder := xml.NewDecoder(body)
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder.Decode(v)
}

func parseTemperature(temp string) float64 {
	var t float64
	fmt.Sscanf(temp, "%f", &t)
	return t
}
