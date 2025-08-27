package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"weather_forecast/internal/models"
)

type mockHTTPClient struct {
	getFunc func(url string) (*http.Response, error)
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	if m.getFunc != nil {
		return m.getFunc(url)
	}
	return nil, nil
}

func createMockResponse(statusCode int, body interface{}) *http.Response {
	jsonData, _ := json.Marshal(body)
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(jsonData)),
		Header:     make(http.Header),
	}
}

func TestGetCityByCEP_ValidCEP(t *testing.T) {
	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			response := models.ViaCepResponse{
				Localidade: "São Paulo",
			}
			return createMockResponse(200, response), nil
		},
	}

	service := &weatherForecast{
		httpClient: mockClient,
		viaCepURL:  "https://viacep.com.br/ws/%s/json/",
		weatherURL: "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=pt",
	}

	city, err := service.GetCityByCEP("01310100")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if city != "São Paulo" {
		t.Errorf("expected 'São Paulo', got %v", city)
	}
}

func TestGetCityByCEP_InvalidCEP(t *testing.T) {
	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			response := models.ViaCepResponse{
				Localidade: "",
				Erro:       "error",
			}
			return createMockResponse(200, response), nil
		},
	}

	service := &weatherForecast{
		httpClient: mockClient,
		viaCepURL:  "https://viacep.com.br/ws/%s/json/",
		weatherURL: "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=pt",
	}

	city, err := service.GetCityByCEP("00000000")

	if err == nil {
		t.Error("expected error for invalid CEP")
	}

	if city != "" {
		t.Errorf("expected empty city, got %v", city)
	}
}

func TestGetCityByCEP_EmptyLocalidade(t *testing.T) {
	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			response := models.ViaCepResponse{
				Localidade: "",
			}
			return createMockResponse(200, response), nil
		},
	}

	service := &weatherForecast{
		httpClient: mockClient,
		viaCepURL:  "https://viacep.com.br/ws/%s/json/",
		weatherURL: "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=pt",
	}

	city, err := service.GetCityByCEP("99999999")

	if err == nil {
		t.Error("expected error for CEP with empty localidade")
	}

	if city != "" {
		t.Errorf("expected empty city, got %v", city)
	}
}

func TestGetCityByCEP_HTTPError(t *testing.T) {
	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			return nil, &mockError{"network error"}
		},
	}

	service := &weatherForecast{
		httpClient: mockClient,
		viaCepURL:  "https://viacep.com.br/ws/%s/json/",
		weatherURL: "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=pt",
	}

	city, err := service.GetCityByCEP("01310100")

	if err == nil {
		t.Error("expected error for HTTP failure")
	}

	if city != "" {
		t.Errorf("expected empty city, got %v", city)
	}
}

func TestGetWeather_Success(t *testing.T) {
	originalKey := os.Getenv("WEATHER_API_KEY")
	os.Setenv("WEATHER_API_KEY", "test-api-key")
	defer os.Setenv("WEATHER_API_KEY", originalKey)

	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			response := models.WeatherAPIResponse{
				Current: struct {
					TempC float64 `json:"temp_c"`
				}{
					TempC: 25.5,
				},
			}
			return createMockResponse(200, response), nil
		},
	}

	service := &weatherForecast{
		httpClient: mockClient,
		viaCepURL:  "https://viacep.com.br/ws/%s/json/",
		weatherURL: "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=pt",
	}

	temp, err := service.GetWeather("São Paulo")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if temp != 25.5 {
		t.Errorf("expected temperature 25.5, got %v", temp)
	}
}

func TestGetWeather_HTTPError(t *testing.T) {
	originalKey := os.Getenv("WEATHER_API_KEY")
	os.Setenv("WEATHER_API_KEY", "test-api-key")
	defer os.Setenv("WEATHER_API_KEY", originalKey)

	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			return nil, &mockError{"weather API error"}
		},
	}

	service := &weatherForecast{
		httpClient: mockClient,
		viaCepURL:  "https://viacep.com.br/ws/%s/json/",
		weatherURL: "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=pt",
	}

	temp, err := service.GetWeather("São Paulo")

	if err == nil {
		t.Error("expected error for HTTP failure")
	}

	if temp != 0 {
		t.Errorf("expected temperature 0, got %v", temp)
	}
}

func TestGetWeather_WithoutAPIKey(t *testing.T) {
	originalKey := os.Getenv("WEATHER_API_KEY")
	os.Setenv("WEATHER_API_KEY", "")
	defer os.Setenv("WEATHER_API_KEY", originalKey)

	mockClient := &mockHTTPClient{
		getFunc: func(url string) (*http.Response, error) {
			return createMockResponse(401, map[string]string{"error": "unauthorized"}), nil
		},
	}

	service := &weatherForecast{
		httpClient: mockClient,
		viaCepURL:  "https://viacep.com.br/ws/%s/json/",
		weatherURL: "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=pt",
	}

	temp, err := service.GetWeather("São Paulo")

	if err != nil {
		t.Errorf("expected no error at service level, got %v", err)
	}

	if temp != 0 {
		t.Errorf("expected temperature 0 for invalid response, got %v", temp)
	}
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}
