package weather

import (
	"errors"
)

// Forecast describes weather conditions for one day in one city
type Forecast struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

// Forecasts define weather conditions for multiple cities
type Forecasts struct {
	Cities map[string]*Forecast `json:"cities"`
}

// NewForecasts return initialized Forecasts
func NewForecasts() *Forecasts {
	return &Forecasts{
		Cities: map[string]*Forecast{},
	}
}

// Errors visible for client

// ErrForecastNotFound means that we couldn't find forecast for given city
var ErrForecastNotFound = errors.New("forecast not found")

// ErrMisconfigured means error in service config
var ErrMisconfigured = errors.New("misconfigured service")

// ErrTooManyRequests means we've exceeded limit in external weather service
var ErrTooManyRequests = errors.New("too many requests")

// ErrInternal means any other error
var ErrInternal = errors.New("internal error")
