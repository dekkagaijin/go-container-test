# Go Weather API Server

A modern Golang web server built with Chi router that provides weather information based on US zip codes.

## Features

- **GET /weather**: Returns current weather data for a given zip code
- **GET /health**: Health check endpoint
- **GET /**: API documentation and usage instructions
- **API Versioning**: `/api/v1/` endpoints for future compatibility
- **Chi Router**: Lightweight, fast HTTP router with middleware support
- **Middleware**: Request logging, panic recovery, CORS, and JSON headers
- Supports major US zip codes
- Returns weather data in JSON format
- Works with OpenWeatherMap API or provides demo data

## Quick Start

1. **Install dependencies:**

   ```bash
   go mod tidy
   ```

2. **Run the server:**

   ```bash
   go run main.go
   ```

3. **Test the API:**

   ```bash
   # Get weather for New York (zip code 10001)
   curl "http://localhost:8080/weather?zip_code=10001"
   
   # Using versioned API
   curl "http://localhost:8080/api/v1/weather?zip_code=10001"
   
   # Check server health
   curl http://localhost:8080/health
   
   # View API documentation
   curl http://localhost:8080/
   ```

## API Endpoints

### Core Endpoints

#### GET /weather?zip_code=XXXXX

#### GET /api/v1/weather?zip_code=XXXXX

Returns weather information for the specified zip code.

**Parameters:**

- `zip_code` (required): 5-digit US zip code (format: XXXXX or XXXXX-XXXX)

**Response:**

```json
{
  "zip_code": "10001",
  "location": "New York",
  "temperature": 72.5,
  "description": "partly cloudy",
  "humidity": 65,
  "wind_speed": 8.2
}
```

#### GET /health

#### GET /api/v1/health

Returns server health status.

**Response:**

```json
{
  "status": "healthy",
  "service": "weather-api"
}
```

#### GET /

Returns API documentation and available endpoints.

### API Versioning

The server supports both unversioned and versioned endpoints:

- **Current**: `/weather`, `/health`
- **Versioned**: `/api/v1/weather`, `/api/v1/health`

Use versioned endpoints for production applications to ensure compatibility with future updates.

**Sample Zip Codes:**

- 10001: New York, NY
- 90210: Beverly Hills, CA
- 60601: Chicago, IL
- 94102: San Francisco, CA
- 77001: Houston, TX
- 33101: Miami, FL
- 98101: Seattle, WA
- 02101: Boston, MA
- 30301: Atlanta, GA
- 75201: Dallas, TX
- 20001: Washington, DC
- 89101: Las Vegas, NV
- 80201: Denver, CO
- 85001: Phoenix, AZ
- 19101: Philadelphia, PA

## Architecture

### Chi Router Features

- **Lightweight**: Minimal overhead HTTP router
- **Middleware Support**: Built-in and custom middleware
- **Route Groups**: Clean API versioning with route groups
- **Method Routing**: Explicit HTTP method handling

### Middleware Stack

1. **Logger**: Logs all HTTP requests with timing
2. **Recoverer**: Gracefully handles panics without crashing
3. **RequestID**: Adds unique request IDs for tracing
4. **RealIP**: Extracts real client IP from headers
5. **JSON/CORS**: Sets appropriate headers for JSON APIs

## Configuration

### Dependencies

- `github.com/go-chi/chi/v5`: HTTP router and middleware

### Environment Variables

- `PORT`: Server port (default: 8080)
- `OPENWEATHER_API_KEY`: OpenWeatherMap API key (optional)

### Using Real Weather Data

To get live weather data instead of demo data:

1. Sign up for a free API key at [OpenWeatherMap](https://openweathermap.org/api)
2. Set the environment variable:

   ```bash
   export OPENWEATHER_API_KEY=your_api_key_here
   go run main.go
   ```

Without an API key, the server returns realistic demo data for testing purposes.

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- `400 Bad Request`: Missing or invalid zip code
- `404 Not Found`: Unsupported zip code or route
- `405 Method Not Allowed`: Unsupported HTTP methods
- `500 Internal Server Error`: Server or external API errors

## Example Usage

```bash
# Valid requests
curl "http://localhost:8080/weather?zip_code=10001"
curl "http://localhost:8080/api/v1/weather?zip_code=90210"

# Health checks
curl "http://localhost:8080/health"
curl "http://localhost:8080/api/v1/health"

# Invalid zip code format
curl "http://localhost:8080/weather?zip_code=123"
# Returns: {"error":"zip_code must be in format XXXXX or XXXXX-XXXX"}

# Missing parameter
curl "http://localhost:8080/weather"
# Returns: {"error":"zip_code parameter is required"}

# Extended zip code format
curl "http://localhost:8080/weather?zip_code=10001-1234"
# Returns: {"zip_code":"10001-1234","location":"New York","temperature":72.5,...}
```

## Development

### Running

```bash
# Development mode
go run main.go

# Build binary
go build -o weather-server main.go
./weather-server
```

### Features Added with Chi

- **Better Routing**: More flexible and performant than net/http
- **Middleware Pipeline**: Composable middleware for cross-cutting concerns
- **Route Groups**: Clean API versioning and organization
- **Request Context**: Enhanced request context with middleware data
- **CORS Support**: Built-in CORS handling for web applications
