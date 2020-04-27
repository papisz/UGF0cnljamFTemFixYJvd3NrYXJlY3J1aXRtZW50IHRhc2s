package weathersrc

import (
	"errors"
	"fmt"
	"log"

	"github.com/papisz/weather"
)

type ForecastManager interface {
	GetForecasts(cities ...string) (*weather.Forecasts, error)
}

type ForecastManagerImpl struct {
	externalProvider ForecastProvider
	storageProvider  WriteableForecastProvider
}

type Option func(o *ForecastManagerImpl)

func NewForecastManager(opts ...Option) *ForecastManagerImpl {
	manager := &ForecastManagerImpl{}
	for _, o := range opts {
		o(manager)
	}
	return manager
}

func WithExternalProvider(provider ForecastProvider) Option {
	return func(m *ForecastManagerImpl) {
		m.externalProvider = provider
	}
}

func WithStorageProvider(provider WriteableForecastProvider) Option {
	return func(m *ForecastManagerImpl) {
		m.storageProvider = provider
	}
}

type ForecastProvider interface {
	GetForecast(city string) (*weather.Forecast, error)
}

type WriteableForecastProvider interface {
	ForecastProvider
	SaveForecast(city string, forecast *weather.Forecast) error
}

// GetForecasts returns forecasts for list of cities
func (m *ForecastManagerImpl) GetForecasts(cities ...string) (*weather.Forecasts, error) {
	forecasts := weather.NewForecasts()

	for _, city := range cities {
		var forecast *weather.Forecast
		var err error

		if forecast, err = m.storageProvider.GetForecast(city); err != nil {
			if !errors.Is(err, weather.ErrForecastNotFound) {
				return nil, fmt.Errorf("error fetching forecast from storage for %s: %w", city, err)
			}

			log.Printf("cache miss for %s", city)

			if forecast, err = m.externalProvider.GetForecast(city); err != nil {
				return nil, fmt.Errorf("error fetching forecast from external provider for %s: %w", city, err)
			}

			if err := m.storageProvider.SaveForecast(city, forecast); err != nil {
				return nil, fmt.Errorf("error saving forecast for %s: %w", city, err)
			}
		} else {
			log.Printf("cache hit for %s", city)
		}
		forecasts.Cities[city] = forecast

	}
	return forecasts, nil
}
