package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// WeatherResponse represents the structure of weather data we'll return
type WeatherResponse struct {
	ZipCode     string  `json:"zip_code"`
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
}

// OpenWeatherMap API response structure (simplified)
type OpenWeatherAPIResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

// ZipCodeLocation maps zip codes to cities (sample mapping)
var zipCodeToCity = map[string]string{
	"10001": "New York,NY,US",
	"90210": "Beverly Hills,CA,US",
	"60601": "Chicago,IL,US",
	"94102": "San Francisco,CA,US",
	"77001": "Houston,TX,US",
	"33101": "Miami,FL,US",
	"98101": "Seattle,WA,US",
	"02101": "Boston,MA,US",
	"30301": "Atlanta,GA,US",
	"75201": "Dallas,TX,US",
	"20001": "Washington,DC,US",
	"89101": "Las Vegas,NV,US",
	"80201": "Denver,CO,US",
	"85001": "Phoenix,AZ,US",
	"19101": "Philadelphia,PA,US",
}

func getWeatherByZipCode(zipCode string) (*WeatherResponse, error) {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		// For demo purposes, return mock data if no API key is provided
		city, exists := zipCodeToCity[zipCode]
		location := "Unknown Location"
		if exists {
			location = strings.Split(city, ",")[0]
		}
		return &WeatherResponse{
			ZipCode:     zipCode,
			Location:    location,
			Temperature: 72.5,
			Description: "partly cloudy (demo data)",
			Humidity:    65,
			WindSpeed:   8.2,
		}, nil
	}

	// Build API URL - OpenWeatherMap supports zip code directly
	baseURL := "http://api.openweathermap.org/data/2.5/weather"
	params := url.Values{}
	params.Add("zip", zipCode+",US") // Assuming US zip codes
	params.Add("appid", apiKey)
	params.Add("units", "imperial") // Fahrenheit

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make HTTP request
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse JSON response
	var apiResp OpenWeatherAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %v", err)
	}

	// Convert to our response format
	description := "clear"
	if len(apiResp.Weather) > 0 {
		description = apiResp.Weather[0].Description
	}

	return &WeatherResponse{
		ZipCode:     zipCode,
		Location:    apiResp.Name,
		Temperature: apiResp.Main.Temp,
		Description: description,
		Humidity:    apiResp.Main.Humidity,
		WindSpeed:   apiResp.Wind.Speed,
	}, nil
}

// Middleware to set JSON content type and CORS headers
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Weather handler using Chi
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	// Get zip code from query parameter
	zipCode := r.URL.Query().Get("zip_code")
	if zipCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "zip_code parameter is required"})
		return
	}

	// Validate zip code format (5 digits, optionally followed by -4 digits)
	zipRegex := regexp.MustCompile(`^\d{5}(-\d{4})?$`)
	if !zipRegex.MatchString(zipCode) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "zip_code must be in format XXXXX or XXXXX-XXXX"})
		return
	}

	// Get weather data
	weather, err := getWeatherByZipCode(zipCode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Return weather data as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "healthy",
		"service": "weather-api",
	}
	json.NewEncoder(w).Encode(response)
}

// Root handler with API documentation
func rootHandler(w http.ResponseWriter, r *http.Request) {
	usage := map[string]interface{}{
		"service": "Weather API Server",
		"endpoints": map[string]string{
			"GET /weather?zip_code=XXXXX": "Get weather by zip code (5 digits)",
			"GET /health":                 "Health check endpoint",
		},
		"example":             "GET /weather?zip_code=10001",
		"supported_zip_codes": []string{"10001", "90210", "60601", "94102", "77001", "33101", "98101", "02101", "30301", "75201", "20001", "89101", "80201", "85001", "19101"},
	}
	json.NewEncoder(w).Encode(usage)
}

func main() {
	// Create Chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)    // Log API request details
	r.Use(middleware.Recoverer) // Recover from panics without crashing server
	r.Use(middleware.RequestID) // Add request ID to context
	r.Use(middleware.RealIP)    // Set RemoteAddr to real client IP
	r.Use(jsonMiddleware)       // Set JSON headers and CORS

	// Define routes
	r.Get("/", rootHandler)
	r.Get("/health", healthHandler)
	r.Get("/weather", weatherHandler)

	// API versioning route group (optional)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/weather", weatherHandler)
		r.Get("/health", healthHandler)
	})

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting weather server with Chi router on port %s...\n", port)
	fmt.Printf("Endpoints available:\n")
	fmt.Printf("  GET /weather?zip_code=10001\n")
	fmt.Printf("  GET /health\n")
	fmt.Printf("  GET /api/v1/weather?zip_code=10001\n")
	fmt.Printf("  GET /api/v1/health\n")

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
