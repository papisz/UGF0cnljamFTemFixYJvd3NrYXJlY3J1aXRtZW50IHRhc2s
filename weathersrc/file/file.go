package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/papisz/weather"
)

// FileWeatherSrc is a weather source just for testing purposes.
type FileWeatherSrc struct {
	path string
}

type Option func(provider *FileWeatherSrc)

func NewWeatherSrc(opts ...Option) *FileWeatherSrc {
	provider := &FileWeatherSrc{}

	for _, opt := range opts {
		opt(provider)
	}

	return provider
}

// WithDirPath sets base path for weather files
func WithDirPath(path string) Option {
	return func(provider *FileWeatherSrc) {
		provider.path = path
	}
}

func (p *FileWeatherSrc) GetForecast(city string) (*weather.Forecast, error) {
	jsonFile, err := os.Open(path.Join(p.path, strings.ToLower(city)+".json"))
	if err != nil {
		return nil, weather.ErrForecastNotFound
	}

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %v", err)
	}

	forecast := &weather.Forecast{}
	if err := json.Unmarshal(jsonBytes, forecast); err != nil {
		return nil, fmt.Errorf("unable to unmarshal bytes: %v", err)
	}
	return forecast, nil
}
