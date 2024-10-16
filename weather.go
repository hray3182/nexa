package main

import (
	"encoding/json"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type WeatherData struct {
	Temperature float64 `json:"temperature"`
	Unit        string  `json:"unit"`
	Description string  `json:"description"`
}

func defineWeatherTool() openai.Tool {
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"location": {
				Type:        jsonschema.String,
				Description: "城市和州，例如 San Francisco, CA",
			},
			"unit": {
				Type: jsonschema.String,
				Enum: []string{"celsius", "fahrenheit"},
			},
		},
		Required: []string{"location"},
	}
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "get_current_weather",
			Description: "獲取指定位置的當前天氣",
			Parameters:  params,
		},
	}
}

func getCurrentWeather(location string, unit string) WeatherData {
	weatherMap := map[string]WeatherData{
		"Boston":            {Temperature: 72, Unit: "fahrenheit", Description: "Sunny"},
		"Boston, MA":        {Temperature: 72, Unit: "fahrenheit", Description: "Sunny"},
		"New York":          {Temperature: 22, Unit: "celsius", Description: "Cloudy"},
		"New York, NY":      {Temperature: 22, Unit: "celsius", Description: "Cloudy"},
		"San Francisco":     {Temperature: 68, Unit: "fahrenheit", Description: "Foggy"},
		"San Francisco, CA": {Temperature: 68, Unit: "fahrenheit", Description: "Foggy"},
	}

	weather, exists := weatherMap[location]
	if !exists {
		return WeatherData{Temperature: 0, Unit: unit, Description: "Unknown location"}
	}

	if unit != "" && unit != weather.Unit {
		if unit == "celsius" && weather.Unit == "fahrenheit" {
			weather.Temperature = (weather.Temperature - 32) * 5 / 9
		} else if unit == "fahrenheit" && weather.Unit == "celsius" {
			weather.Temperature = weather.Temperature*9/5 + 32
		}
		weather.Unit = unit
	}

	return weather
}

func handleWeatherCall(call openai.ToolCall) openai.ChatCompletionMessage {
	var params struct {
		Location string `json:"location"`
		Unit     string `json:"unit"`
	}
	json.Unmarshal([]byte(call.Function.Arguments), &params)
	weatherData := getCurrentWeather(params.Location, params.Unit)
	weatherJSON, _ := json.Marshal(weatherData)
	return openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    string(weatherJSON),
		Name:       call.Function.Name,
		ToolCallID: call.ID,
	}
}
