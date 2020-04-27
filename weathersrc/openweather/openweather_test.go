package openweather

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/papisz/weather"
	"github.com/papisz/weather/testutils"
)

func TestOpenWeatherSrc_GetForecast(t *testing.T) {
	type fields struct {
		apiKey string
		status int
		body   []byte
	}
	type args struct {
		city string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        *weather.Forecast
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Proper response with forecast",
			fields: fields{
				apiKey: "fake",
				status: http.StatusOK,
				body:   testutils.JSONFileToBytes("../../testdata/source", "london.json"),
			},
			args:    args{city: "London"},
			want:    testutils.ForecastFromJSON("london.json"),
			wantErr: false,
		},
		{
			name: "City couldn't be found",
			fields: fields{
				apiKey: "fake",
				status: http.StatusNotFound,
				body:   []byte{},
			},
			args:        args{city: "London"},
			want:        nil,
			wantErr:     true,
			expectedErr: weather.ErrForecastNotFound,
		},
		{
			name: "Unauthorized request",
			fields: fields{
				apiKey: "fake",
				status: http.StatusUnauthorized,
				body:   []byte{},
			},
			args:        args{city: "London"},
			want:        nil,
			wantErr:     true,
			expectedErr: weather.ErrMisconfigured,
		},
		{
			name: "Request limit exceeded",
			fields: fields{
				apiKey: "fake",
				status: http.StatusTooManyRequests,
				body:   []byte{},
			},
			args:        args{city: "London"},
			want:        nil,
			wantErr:     true,
			expectedErr: weather.ErrTooManyRequests,
		},
		{
			name: "Unknown error",
			fields: fields{
				apiKey: "fake",
				status: http.StatusTeapot,
				body:   []byte{},
			},
			args:        args{city: "London"},
			want:        nil,
			wantErr:     true,
			expectedErr: errors.New("external service returned 418 and "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := func(status int, body []byte) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Query().Get("appid"), tt.fields.apiKey)
					assert.Equal(t, req.URL.Query().Get("q"), tt.args.city)
					res.WriteHeader(status)
					res.Write([]byte(body))
				}))
			}(tt.fields.status, tt.fields.body)
			p := NewWeatherSrc(
				WithURL(testServer.URL),
				WithDefaultClient(),
				WithAPIKey(tt.fields.apiKey),
			)
			got, err := p.GetForecast(tt.args.city)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenWeatherSrc.GetForecast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenWeatherSrc.GetForecast() = %v, want %v", got, tt.want)
			}

			switch tt.wantErr {
			case true:
				assert.EqualError(t, tt.expectedErr, err.Error())
			case false:
				assert.Nil(t, err)
			}

			testServer.Close()
		})
	}
}
