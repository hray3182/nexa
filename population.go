package main

import (
	"encoding/json"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func definePopulationTool() openai.Tool {
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"city": {
				Type:        jsonschema.String,
				Description: "城市名稱",
			},
		},
		Required: []string{"city"},
	}
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "get_population",
			Description: "獲取指定城市的人口數量",
			Parameters:  params,
		},
	}
}

func getPopulation(city string) int {
	populationMap := map[string]int{
		"Boston":        675647,
		"New York":      8336817,
		"San Francisco": 873965,
	}
	return populationMap[city]
}

func handlePopulationCall(call openai.ToolCall) openai.ChatCompletionMessage {
	var params struct {
		City string `json:"city"`
	}
	json.Unmarshal([]byte(call.Function.Arguments), &params)
	population := getPopulation(params.City)
	populationJSON, _ := json.Marshal(map[string]int{"population": population})
	return openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    string(populationJSON),
		Name:       call.Function.Name,
		ToolCallID: call.ID,
	}
}
