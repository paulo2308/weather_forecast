package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"weather_forecast/internal/models"
)

type mockWeatherService struct {
	getCityByCEPFunc func(cep string) (string, error)
	getWeatherFunc   func(city string) (float64, error)
}

func (m *mockWeatherService) GetCityByCEP(cep string) (string, error) {
	if m.getCityByCEPFunc != nil {
		return m.getCityByCEPFunc(cep)
	}
	return "", nil
}

func (m *mockWeatherService) GetWeather(city string) (float64, error) {
	if m.getWeatherFunc != nil {
		return m.getWeatherFunc(city)
	}
	return 0, nil
}

func TestInvalidCEP(t *testing.T) {
	mockService := &mockWeatherService{}
	handler := NewWeatherHandler(mockService)

	req, _ := http.NewRequest("GET", "/weather?cep=123", nil)
	rr := httptest.NewRecorder()

	handler.WeatherHandler(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 got %v", status)
	}

	var response models.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response.Message != "invalid zipcode" {
		t.Errorf("expected 'invalid zipcode' got %v", response.Message)
	}
}

func TestCEPNotFound(t *testing.T) {
	mockService := &mockWeatherService{
		getCityByCEPFunc: func(cep string) (string, error) {
			return "", &mockError{"not found"}
		},
	}
	handler := NewWeatherHandler(mockService)

	req, _ := http.NewRequest("GET", "/weather?cep=12345678", nil)
	rr := httptest.NewRecorder()

	handler.WeatherHandler(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("expected 404 got %v", status)
	}

	var response models.ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response.Message != "can not find zipcode" {
		t.Errorf("expected 'can not find zipcode' got %v", response.Message)
	}
}

func TestWeatherAPIError(t *testing.T) {
	mockService := &mockWeatherService{
		getCityByCEPFunc: func(cep string) (string, error) {
			return "São Paulo", nil
		},
		getWeatherFunc: func(city string) (float64, error) {
			return 0, &mockError{"weather API error"}
		},
	}
	handler := NewWeatherHandler(mockService)

	req, _ := http.NewRequest("GET", "/weather?cep=01310100", nil)
	rr := httptest.NewRecorder()

	handler.WeatherHandler(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected 500 got %v", status)
	}
}

func TestCEPOk(t *testing.T) {
	mockService := &mockWeatherService{
		getCityByCEPFunc: func(cep string) (string, error) {
			return "São Paulo", nil
		},
		getWeatherFunc: func(city string) (float64, error) {
			return 25.5, nil
		},
	}
	handler := NewWeatherHandler(mockService)

	req, _ := http.NewRequest("GET", "/weather?cep=01310100", nil)
	rr := httptest.NewRecorder()

	handler.WeatherHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected 200 got %v", status)
	}

	var response models.WeatherResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	expectedTempC := 25.5
	expectedTempF := 25.5*1.8 + 32
	expectedTempK := 25.5 + 273

	if response.TempC != expectedTempC {
		t.Errorf("expected temp_C %v got %v", expectedTempC, response.TempC)
	}
	if response.TempF != expectedTempF {
		t.Errorf("expected temp_F %v got %v", expectedTempF, response.TempF)
	}
	if response.TempK != expectedTempK {
		t.Errorf("expected temp_K %v got %v", expectedTempK, response.TempK)
	}
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}
