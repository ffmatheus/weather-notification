package service

import (
	"context"
	"time"
	"weather-notification/internal/domain/entity"
	handler "weather-notification/internal/domain/error_handler"
	"weather-notification/internal/domain/repository"

	"github.com/google/uuid"
)

type CPTECClient interface {
	SearchCities(ctx context.Context, cityName string) ([]entity.Location, error)
	GetWeatherForecast(ctx context.Context, cptecCode int) (*entity.WeatherForecastCollection, error)
	GetWaveForecast(ctx context.Context, cptecCode int, date time.Time) (*entity.WaveInfo, error)
}

type WeatherService struct {
	cptecClient  CPTECClient
	locationRepo repository.LocationRepository
}

func NewWeatherService(client CPTECClient, locationRepo repository.LocationRepository) *WeatherService {
	return &WeatherService{
		cptecClient:  client,
		locationRepo: locationRepo,
	}
}

func (s *WeatherService) SearchLocation(ctx context.Context, cityName string) ([]entity.Location, error) {
	location, err := s.locationRepo.FindByNameAndState(ctx, cityName, "")
	if err == nil {
		return []entity.Location{*location}, nil
	}

	locations, err := s.cptecClient.SearchCities(ctx, cityName)
	if err != nil {
		return nil, err
	}

	var newLocations []entity.Location

	for _, loc := range locations {
		location, err := s.locationRepo.FindByCPTECCode(ctx, loc.CPTECCode)
		if err == nil {
			newLocations = append(newLocations, *location)
		} else if err == handler.ErrNotFound {
			_ = s.locationRepo.Create(ctx, &loc)
			newLocations = append(newLocations, loc)
		}
	}

	return newLocations, nil
}

func (s *WeatherService) GetForecast(ctx context.Context, locationID uuid.UUID) (*entity.WeatherForecastCollection, error) {
	location, err := s.locationRepo.FindByID(ctx, locationID)
	if err != nil {
		return nil, err
	}

	forecast, err := s.cptecClient.GetWeatherForecast(ctx, location.CPTECCode)
	if err != nil {
		return nil, err
	}

	for i, f := range forecast.Forecasts {
		forecastDate := f.Date

		wave, err := s.cptecClient.GetWaveForecast(ctx, location.CPTECCode, forecastDate)
		if err == nil {
			forecast.Forecasts[i].Wave = wave
		}
	}

	return forecast, nil
}
