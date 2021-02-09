package data

import (
	"github.com/jadilet/weather/api"
)

// GetWeatherByCity get weather by city from an API
func GetWeatherByCity(city string) (*api.QueryResponse, error) {

	mw := &api.MetaWeather{}

	var qr *api.QueryResponse
	var err error

	qr, err = mw.GetWeather(city)

	if err != nil {
		return qr, err
	}

	return qr, nil
}
