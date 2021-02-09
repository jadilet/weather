package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jadilet/weather/data"
)

// Weathers controller
type Weathers struct {
	l    *log.Logger
	tmpl *template.Template
}

// NewWeathers a new instance of Weathers
func NewWeathers(l *log.Logger, t *template.Template) *Weathers {
	return &Weathers{l, t}
}

// GetIndexPage render index page
func (weather *Weathers) GetIndexPage(w http.ResponseWriter, r *http.Request) {
	weather.l.Println("handle GET weather request")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var city string

	city = r.FormValue("search")

	if city == "" {
		city = "warsaw" // default city name
	}

	response, err := data.GetWeatherByCity(city)

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		err := weather.tmpl.Execute(w, response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
	}
}
