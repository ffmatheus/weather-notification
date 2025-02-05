package handler

import "time"

//WEATHER

type SearchLocationResponse struct {
	ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name  string `json:"name" example:"São Paulo"`
	State string `json:"state" example:"SP"`
}

type ForecastResponse struct {
	Date     time.Time `json:"date" example:"2024-02-02T00:00:00Z"`
	MinTemp  float64   `json:"min_temp" example:"18.5"`
	MaxTemp  float64   `json:"max_temp" example:"27.8"`
	Forecast string    `json:"forecast" example:"Parcialmente nublado"`
	IUV      float64   `json:"uv" example:"3.5"`
	Wave     *WaveInfo `json:"wave,omitempty"`
}

type WaveInfo struct {
	Height    float64 `json:"height" example:"1.5"`
	Direction string  `json:"direction" example:"Sudeste"`
}

//USER

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	City  string `json:"city" binding:"required" example:"São Paulo"`
}

type UpdateUserRequest struct {
	Name string `json:"name,omitempty"`
	City string `json:"city,omitempty"`
}

type ToggleOptOutRequest struct {
	OptOut bool `json:"opt_out" binding:"boolean"`
}
