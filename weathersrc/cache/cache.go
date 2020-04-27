package cache

import (
	"time"

	"github.com/papisz/weather"
	gocache "github.com/patrickmn/go-cache"
)

type CacheWeatherSrc struct {
	cache *gocache.Cache
	ttl   time.Duration
}

type Option func(provider *CacheWeatherSrc)

func NewWeatherSrc(opts ...Option) *CacheWeatherSrc {
	provider := &CacheWeatherSrc{}

	for _, opt := range opts {
		opt(provider)
	}

	return provider
}

func WithTTL(ttl time.Duration) Option {
	return func(provider *CacheWeatherSrc) {
		provider.ttl = ttl
		provider.cache = gocache.New(provider.ttl, provider.ttl)
	}
}

func (p *CacheWeatherSrc) GetForecast(city string) (*weather.Forecast, error) {
	if x, found := p.cache.Get(city); found {
		return x.(*weather.Forecast), nil
	}

	return nil, weather.ErrForecastNotFound
}

func (p *CacheWeatherSrc) SaveForecast(city string, forecast *weather.Forecast) error {
	p.cache.Set(city, forecast, gocache.DefaultExpiration)
	return nil
}
