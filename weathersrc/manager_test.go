package weathersrc

import (
	"testing"

	"github.com/papisz/weather"
	"github.com/stretchr/testify/mock"
)

func TestForecastManagerImpl_GetForecasts(t *testing.T) {

	tests := []struct {
		name             string
		storageProvider  *returnedForecast
		externalProvider *returnedForecast
		wantErr          bool
	}{
		{
			name: "Storage miss, external provider hit",
			storageProvider: &returnedForecast{
				nil,
				weather.ErrForecastNotFound,
			},

			externalProvider: &returnedForecast{
				forecast: &weather.Forecast{},
				err:      nil,
			},

			wantErr: false,
		},
		{
			name: "Storage hit, external provider not called",
			storageProvider: &returnedForecast{
				forecast: &weather.Forecast{},
				err:      nil,
			},

			externalProvider: nil,
			wantErr:          false,
		},
		{
			name: "Storage miss, external provider miss",
			storageProvider: &returnedForecast{
				forecast: nil,
				err:      weather.ErrForecastNotFound,
			},

			externalProvider: &returnedForecast{
				forecast: nil,
				err:      weather.ErrForecastNotFound,
			},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			city := "london"
			externalProvider := newMockProvider(tt.externalProvider)
			storageProvider := newMockProvider(tt.storageProvider)

			m := NewForecastManager(
				WithExternalProvider(externalProvider),
				WithStorageProvider(storageProvider),
			)
			_, err := m.GetForecasts(city)
			if (err != nil) != tt.wantErr {
				t.Errorf("ForecastManagerImpl.GetForecasts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.externalProvider != nil {
				externalProvider.AssertExpectations(t)
			}
			if tt.storageProvider != nil {
				storageProvider.AssertExpectations(t)
			}
		})
	}
}

type MockProvider struct {
	mock.Mock
}

func newMockProvider(returnedForecast *returnedForecast) *MockProvider {
	if returnedForecast == nil {
		return nil
	}

	p := new(MockProvider)
	p.On("GetForecast", mock.Anything).Return(returnedForecast.forecast, returnedForecast.err)
	return p
}

func (m *MockProvider) GetForecast(city string) (*weather.Forecast, error) {
	args := m.Called(city)
	return args.Get(0).(*weather.Forecast), args.Error(1)
}

func (m *MockProvider) SaveForecast(city string, forecast *weather.Forecast) error {
	return nil
}

type returnedForecast struct {
	forecast *weather.Forecast
	err      error
}
