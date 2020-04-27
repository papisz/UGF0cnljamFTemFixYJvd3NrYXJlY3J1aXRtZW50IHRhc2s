package openweather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/papisz/weather"
)

type OpenWeatherSrc struct {
	URL    string
	apiKey string
	client HTTPClient
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Option func(provider *OpenWeatherSrc)

func NewWeatherSrc(opts ...Option) *OpenWeatherSrc {
	provider := &OpenWeatherSrc{}

	for _, opt := range opts {
		opt(provider)
	}

	return provider
}

func WithDefaultClient() Option {
	return func(provider *OpenWeatherSrc) {
		provider.client = &http.Client{Timeout: 10 * time.Second}
	}
}

func WithCustomClient(client HTTPClient) Option {
	return func(provider *OpenWeatherSrc) {
		provider.client = client
	}
}

func WithURL(url string) Option {
	return func(provider *OpenWeatherSrc) {
		provider.URL = url
	}
}

func WithAPIKey(apiKey string) Option {
	return func(provider *OpenWeatherSrc) {
		provider.apiKey = apiKey
	}
}

func (p *OpenWeatherSrc) getURL(city string) string {
	v := url.Values{}
	v.Add("q", city)
	v.Add("appid", p.apiKey)
	return p.URL + "?" + v.Encode()
}

func (p *OpenWeatherSrc) GetForecast(city string) (*weather.Forecast, error) {
	req, err := http.NewRequest(http.MethodGet, p.getURL(city), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, weather.ErrForecastNotFound
	case http.StatusUnauthorized:
		return nil, weather.ErrMisconfigured
	case http.StatusTooManyRequests:
		return nil, weather.ErrTooManyRequests
	default:
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("external service returned %d and %s", resp.StatusCode, body)
	}

	var forecast = weather.Forecast{}
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, err
	}

	return &forecast, nil
}
