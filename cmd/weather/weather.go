package main

import (
	"flag"
	"time"

	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/papisz/weather/api/http"
	"github.com/papisz/weather/weathersrc"
	"github.com/papisz/weather/weathersrc/cache"
	"github.com/papisz/weather/weathersrc/openweather"
	"gopkg.in/relistan/rubberneck.v1"
)

type Config struct {
	Listen           string        `default:"localhost:5555"`
	WeatherSrcAPIKey string        `required:"true"`
	WeatherSrcAPIURL string        `default:"https://api.openweathermap.org/data/2.5/weather"`
	CacheTTL         time.Duration `default:"5h"`
}

func ParseConfig() *Config {
	help := flag.Bool("help", false, "print help")
	flag.Parse()

	var config Config

	if help != nil && *help {
		if err := envconfig.Usage("weather", &config); err != nil {
			log.Fatal(err)
		}
		return nil
	}

	if err := envconfig.Process("weather", &config); err != nil {
		log.Fatalf("invalid config, %v", err)
		envconfig.Usage("weather", &config)
		return nil
	}

	rubberneck.Print(config)
	return &config
}

func main() {
	config := ParseConfig()

	if config == nil {
		return
	}

	if err := http.NewApi(
		http.WithListenAddress(config.Listen),
		http.WithForecastManager(
			weathersrc.NewForecastManager(
				weathersrc.WithExternalProvider(
					openweather.NewWeatherSrc(
						openweather.WithURL(config.WeatherSrcAPIURL),
						openweather.WithAPIKey(config.WeatherSrcAPIKey),
						openweather.WithDefaultClient(),
					),
				),
				weathersrc.WithStorageProvider(
					cache.NewWeatherSrc(
						cache.WithTTL(config.CacheTTL),
					),
				),
			),
		),
	).Serve(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}

}
