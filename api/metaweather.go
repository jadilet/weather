/*
Package api gets data from MetaWeather
Link to meta weather API: https://www.metaweather.com/api/
*/
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	getWeidURL    = "https://www.metaweather.com/api/location/search/"
	getWeatherURL = "https://www.metaweather.com/api/location/"
)

// MetaWeather api
type MetaWeather struct {
}

// Weather defines the structure of an API weather
type Weather struct {
	ID                      float64 `json:"id"`
	WeatherStateName        string  `json:"weather_state_name"`
	WeatherStateAbbr        string  `json:"weather_state_abbr"`
	WeatherDirectionCompass string  `json:"wind_direction_compass"`
	ApplicableDate          string  `json:"applicable_date"`
	MinTemp                 float64 `json:"min_temp"`
	MaxTemp                 float64 `json:"max_temp"`
	Temp                    float64 `json:"the_temp"`
	WindSpeed               float64 `json:"wind_speed"`
	WindDirection           float64 `json:"wind_direction"`
	AirPressure             float64 `json:"air_pressure"`
	Humidity                float64 `json:"humidity"`
	Visibility              float64 `json:"visibility"`
	Predictability          float64 `json:"predictability"`
}

// Search defines the structure of an API
type Search struct {
	Title               string    `json:"title"`
	LocationType        string    `json:"location_type"`
	ConsolidatedWeather []Weather `json:"consolidated_weather"`
}

type query struct {
	Title        string `json:"title"`
	LocationType string `json:"location_type"`
	Woeid        uint64 `json:"woeid"`
	LatLon       string `json:"latt_long"`
}

// WeatherResp dates of weather
type WeatherResp struct {
	WeekdayAbbr string
	TheTemp     string
	IconURL     string
}

// QueryResponse defines structure of json for template/html
type QueryResponse struct {
	Weekday          string
	Date             string
	City             string
	TheTemp          string
	WeatherStateName string
	Humidity         string
	Predictability   string
	WindSpeed        string
	Arr              []*WeatherResp
}

var icons = map[string]string{
	"Snow": "sn",
	"Sleet": "sl",
	"Hail": "h",
	"Thunderstorm": "t",

	"Heavy Rain": "hr",
	"Light Rain": "lr",
	"Showers": "s",
	"Heavy Cloud": "hc",
	"Light Cloud": "lc",
	"Clear": "c",

}

// getweid find id of city or country by query
// By ID fetches the city or countr weathers datum
func (m *MetaWeather) getweid(q string) (*query, error) {
	resp, err := http.Get(fmt.Sprintf("%s?query=%s", getWeidURL, q))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	qr := []query{}

	err = json.NewDecoder(resp.Body).Decode(&qr)

	if err != nil {
		return nil, err
	}

	if len(qr) == 0 {
		return nil, errors.New("Not found query")
	}

	return &qr[0], nil
}

// GetWeather gets data from meta weather API
func (m *MetaWeather) GetWeather(query string) (*QueryResponse, error) {

	qr, err := m.getweid(query)

	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("%s%d", getWeatherURL, qr.Woeid))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	search := Search{}

	err = json.NewDecoder(resp.Body).Decode(&search)

	if err != nil {
		return nil, err
	}

	rp := QueryResponse{}

	for i, v := range search.ConsolidatedWeather {
		t, err := time.Parse("2006-01-02", v.ApplicableDate)

		if i == 0 {
			rp.City = search.Title
			rp.Weekday = t.Weekday().String()
			rp.Humidity = fmt.Sprintf("%.2f", v.Humidity)
			rp.Date = t.Format("January 2, 2006")
			rp.Predictability = fmt.Sprintf("%.2f", v.Predictability)
			rp.WindSpeed = fmt.Sprintf("%.2f", v.WindSpeed)
			rp.TheTemp = fmt.Sprintf("%.2f", v.Temp)

			rp.WeatherStateName = v.WeatherStateName

		}

		if err != nil {
			continue
		}

		res := WeatherResp{}

		res.WeekdayAbbr = t.Weekday().String()[:3]
		if v, ok := icons[v.WeatherStateName]; ok {
			res.IconURL = fmt.Sprintf("https://www.metaweather.com/static/img/weather/%s.svg", v)
		}

		res.TheTemp = fmt.Sprintf("%.2f", v.Temp)
		rp.Arr = append(rp.Arr, &res)
	}

	return &rp, nil
}
