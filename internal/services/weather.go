package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"weather_forecast/internal/models"
)

// HTTPClient interface para permitir mocking
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type WeatherForecast interface {
	GetCityByCEP(cep string) (string, error)
	GetWeather(city string) (float64, error)
}

type weatherForecast struct {
	httpClient HTTPClient
	viaCepURL  string
	weatherURL string
}

func NewWeatherForecast() WeatherForecast {
	return &weatherForecast{
		httpClient: &http.Client{},
		viaCepURL:  "https://viacep.com.br/ws/%s/json/",
		weatherURL: "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=pt",
	}

}

func (wf *weatherForecast) GetCityByCEP(cep string) (string, error) {
	url := fmt.Sprintf(wf.viaCepURL, cep)
	resp, err := wf.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data models.ViaCepResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.Erro == "true" || data.Localidade == "" {
		return "", fmt.Errorf("not found")
	}

	return data.Localidade, nil
}

func (wf *weatherForecast) GetWeather(city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	encodedCity := url.QueryEscape(city)
	weatherURL := fmt.Sprintf(wf.weatherURL, apiKey, encodedCity)

	resp, err := wf.httpClient.Get(weatherURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data models.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
