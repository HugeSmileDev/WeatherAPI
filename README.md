```markdown
# Go Weather Server

This is a simple HTTP server written in Go that fetches weather information from the OpenWeather API based on latitude and longitude coordinates provided by the user. It determines whether it's hot, cold, or moderate outside and returns the weather condition along with the temperature.

## Features

- Exposes an endpoint `/weather` to get weather information based on latitude and longitude coordinates.
- Determines the weather condition (hot, cold, or moderate) based on the temperature.
- Fetches weather data from the OpenWeather API.
- Utilizes environment variables to load the OpenWeather API key securely.

## Requirements

- Go 1.13 or higher installed on your machine.
- OpenWeather API key. You can obtain it by signing up at [OpenWeather](https://home.openweathermap.org/users/sign_up).

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/HugeSmileDev/WeatherAPI.git

2. Navigate to the project directory:

   ```bash
   cd weather-api

3. Create a `.env` file in the project root and add your OpenWeather API key:

   ```bash
   OPENWEATHER_API_KEY=your_openweather_api_key

4. Run the server:

   ```bash
   go run main.go
   ```

## Usage

To get weather information, send a GET request to the `/weather` endpoint with latitude and longitude coordinates as query parameters. For example:

```
http://localhost:8080/weather?lat=25.5922166&lon=85.12761069999999
```

Replace `25.5922166` and `85.12761069999999` with your desired latitude and longitude coordinates.
