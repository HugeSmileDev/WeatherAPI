package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// WeatherResponse struct to unmarshal JSON response from OpenWeather API
type WeatherResponse struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

// GetWeather function fetches weather data from OpenWeather API
func GetWeather(lat, lon float64, apiKey string) (*WeatherResponse, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&APPID=%s", lat, lon, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	weatherResponse := &WeatherResponse{}
	err = json.Unmarshal(body, weatherResponse)
	if err != nil {
		return nil, err
	}

	return weatherResponse, nil
}

// WeatherHandler handles incoming HTTP requests for weather information
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	// Log incoming request
	log.Printf("Incoming request from %s for %s", r.RemoteAddr, r.URL.Path)

	// Parse latitude and longitude from query parameters
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	// Convert latitude and longitude to float64
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		log.Printf("Error parsing latitude: %v", err)
		return
	}
	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		log.Printf("Error parsing longitude: %v", err)
		return
	}

	// Load API key from environment file
	apiKey, err := loadAPIKey()
	if err != nil {
		http.Error(w, "Failed to load API key", http.StatusInternalServerError)
		return
	}

	// Get weather data
	weather, err := GetWeather(latFloat, lonFloat, apiKey)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	// Determine weather condition based on temperature
	var weatherCondition string
	temperature := weather.Main.Temp
	switch {
	case temperature >= 30:
		weatherCondition = "hot"
	case temperature <= 10:
		weatherCondition = "cold"
	default:
		weatherCondition = "moderate"
	}

	// Construct response
	response := fmt.Sprintf("Weather: %s, Temperature: %.1f°C, Condition: %s", weather.Weather[0].Main, temperature, weatherCondition)

	// Log response
	log.Printf("Response sent for request from %s: %s", r.RemoteAddr, response)

	// Send response
	fmt.Fprintf(w, response)
}

// loadAPIKey loads the API key from the environment file
func loadAPIKey() (string, error) {
	file, err := os.Open(".env")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var apiKey string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "OPENWEATHER_API_KEY=") {
			apiKey = strings.TrimPrefix(line, "OPENWEATHER_API_KEY=")
			break
		}
	}
	if apiKey == "" {
		return "", fmt.Errorf("API key not found in .env file")
	}

	return apiKey, nil
}

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Start the HTTP server
	addr := ":8080"
	log.Printf("Server is running on %s", addr)
	http.HandleFunc("/weather", WeatherHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
