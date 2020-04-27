package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/papisz/weather"

	"github.com/papisz/weather/testutils"
	"github.com/papisz/weather/weathersrc"
	"github.com/papisz/weather/weathersrc/cache"
	"github.com/papisz/weather/weathersrc/file"
	"github.com/stretchr/testify/assert"
)

func TestWithFileIntegration(t *testing.T) {
	r := requestCreator{listenAddress: "localhost:5555"}

	type args struct {
		url string
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name: "Successfully get forecasts for two cities",
			args: args{
				url: "forecast?city=london&city=warsaw",
			},
			expectedStatus: 200,
			expectedBody:   testutils.JSONFileToBytes("../../testdata/response/", "1.json"),
		},
		{
			name: "Error: get forecasts for two cities - one is missing",
			args: args{
				url: "forecast?city=london&city=szczebrzeszyn",
			},
			expectedStatus: 404,
			expectedBody: []byte(`{
				"status": "forecast not found"
			}`),
		},
		{
			name: "Error: no cities given",
			args: args{
				url: "forecast",
			},
			expectedStatus: 400,
			expectedBody: []byte(`{
				"status": "unable to parse cities"
			}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewApi(
				WithListenAddress(r.listenAddress),
				WithForecastManager(weathersrc.NewForecastManager(
					weathersrc.WithExternalProvider(
						file.NewWeatherSrc(
							file.WithDirPath("../../testdata/source"),
						),
					),
					weathersrc.WithStorageProvider(
						cache.NewWeatherSrc(cache.WithTTL(5*time.Second)),
					),
				)),
			)
			w := httptest.NewRecorder()
			a.GetForecasts(w, r.newRequest(tt.args.url))

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, string(tt.expectedBody), w.Body.String())
		})

	}
}

func TestHTTPApi_GetForecastsErrors(t *testing.T) {
	r := requestCreator{listenAddress: "localhost:5555"}

	type fields struct {
		weatherManager weathersrc.ForecastManager
	}

	tests := []struct {
		name           string
		fields         fields
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name: "Get forecasts returning misconfigured service",
			fields: fields{
				weatherManager: &forecastManagerMock{
					err: fmt.Errorf("%w", weather.ErrMisconfigured),
				},
			},
			expectedStatus: 500,
			expectedBody: []byte(`{
				"status": "misconfigured service"
			}`),
		},
		{
			name: "Get forecasts with too many requests",
			fields: fields{
				weatherManager: &forecastManagerMock{
					err: fmt.Errorf("%w", weather.ErrTooManyRequests),
				},
			},
			expectedStatus: 500,
			expectedBody: []byte(`{
				"status": "too many requests"
			}`),
		},
		{
			name: "Get forecasts with other error",

			fields: fields{
				weatherManager: &forecastManagerMock{
					err: fmt.Errorf("other error"),
				},
			},
			expectedStatus: 500,
			expectedBody: []byte(`{
				"status": "internal error"
			}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewApi(
				WithListenAddress(r.listenAddress),
				WithForecastManager(tt.fields.weatherManager),
			)
			w := httptest.NewRecorder()
			a.GetForecasts(w, r.newRequest("forecast?city=london&city=szczebrzeszyn"))

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, string(tt.expectedBody), w.Body.String())
		})

	}
}

type requestCreator struct {
	listenAddress string
}

func (r *requestCreator) newRequest(url string) *http.Request {
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", r.listenAddress, url), nil)
	return request
}

type forecastManagerMock struct {
	err       error
	forecasts *weather.Forecasts
}

func (m *forecastManagerMock) GetForecasts(cities ...string) (*weather.Forecasts, error) {
	return m.forecasts, m.err
}
