package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"weather_forecast/internal/models"
	"weather_forecast/internal/services"
)

var cepRegex = regexp.MustCompile(`^\d{8}$`)

type WeatherHandler interface {
	WeatherHandler(w http.ResponseWriter, r *http.Request)
}
type weatherForecast struct {
	swf services.WeatherForecast
}

func NewWeatherHandler(swf services.WeatherForecast) WeatherHandler {
	return &weatherForecast{
		swf: swf,
	}
}

func (wf *weatherForecast) WeatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")

	if !cepRegex.MatchString(cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "invalid zipcode"})
		return
	}

	city, err := wf.swf.GetCityByCEP(cep)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "can not find zipcode"})
		return
	}

	tempC, err := wf.swf.GetWeather(city)
	if err != nil {
		http.Error(w, "failed to fetch weather", http.StatusInternalServerError)
		return
	}

	resp := models.WeatherResponse{
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
