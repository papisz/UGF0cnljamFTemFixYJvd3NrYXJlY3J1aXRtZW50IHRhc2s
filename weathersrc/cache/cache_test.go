package cache

import (
	"testing"
	"time"

	"github.com/papisz/weather"
	"github.com/papisz/weather/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCacheWeatherSrc_GetSaveForecast(t *testing.T) {
	type args struct {
		city     string
		forecast *weather.Forecast
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "#1 City",
			args: args{
				city:     "London",
				forecast: testutils.ForecastFromJSON("london.json"),
			},
		},
		{
			name: "#1 City",
			args: args{
				city:     "Warsaw",
				forecast: testutils.ForecastFromJSON("warsaw.json"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewWeatherSrc(WithTTL(5 * time.Second))
			forecast, err := p.GetForecast(tt.args.city)

			assert.Equal(t, weather.ErrForecastNotFound, err)
			assert.Nil(t, forecast)

			if err := p.SaveForecast(tt.args.city, tt.args.forecast); (err != nil) != tt.wantErr {
				t.Errorf("CacheWeatherSrc.SaveForecast() error = %v, wantErr %v", err, tt.wantErr)
			}

			forecast, err = p.GetForecast(tt.args.city)

			assert.Nil(t, err)
			assert.Equal(t, tt.args.forecast, forecast)
		})
	}
}
