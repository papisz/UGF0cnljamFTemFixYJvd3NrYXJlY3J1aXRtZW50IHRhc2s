package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/papisz/weather"
	"github.com/papisz/weather/weathersrc"
)

type HTTPApi struct {
	ListenAddress  string
	WeatherManager weathersrc.ForecastManager
}

type Option func(api *HTTPApi)

func NewApi(opts ...Option) *HTTPApi {
	api := &HTTPApi{}

	for _, opt := range opts {
		opt(api)
	}
	return api
}

func WithListenAddress(address string) Option {
	return func(api *HTTPApi) {
		api.ListenAddress = address
	}
}

func WithForecastManager(m weathersrc.ForecastManager) Option {
	return func(api *HTTPApi) {
		api.WeatherManager = m
	}
}

func (a *HTTPApi) GetForecasts(w http.ResponseWriter, r *http.Request) {
	var err error
	var forecasts = weather.NewForecasts()
	var cities []string

	parseCities := func(r *http.Request) []string {
		cities, ok := r.URL.Query()["city"]
		if !ok {
			return nil
		}
		return cities
	}

	if cities = parseCities(r); cities == nil {
		render.Render(w, r, &ErrResponse{
			Err:            nil,
			HTTPStatusCode: http.StatusBadRequest,
			StatusText:     "unable to parse cities",
		})
		return
	}

	if forecasts, err = a.WeatherManager.GetForecasts(cities...); err != nil {
		switch errors.Unwrap(err) {
		case weather.ErrForecastNotFound:
			render.Render(w, r, &ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusNotFound,
				StatusText:     weather.ErrForecastNotFound.Error(),
			})
		case weather.ErrMisconfigured:
			render.Render(w, r, &ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				StatusText:     weather.ErrMisconfigured.Error(),
			})
		case weather.ErrTooManyRequests:
			render.Render(w, r, &ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				StatusText:     weather.ErrTooManyRequests.Error(),
			})
		default:
			render.Render(w, r, &ErrResponse{
				Err:            err,
				HTTPStatusCode: http.StatusInternalServerError,
				StatusText:     weather.ErrInternal.Error(),
			})
		}
		return
	}

	render.JSON(w, r, forecasts)
	return
}

func (a *HTTPApi) Serve() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/forecast", a.GetForecasts)
	return http.ListenAndServe(a.ListenAddress, r)
}
