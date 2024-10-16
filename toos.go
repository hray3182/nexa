package main

import "github.com/sashabaranov/go-openai"

func defineTools() []openai.Tool {
	return []openai.Tool{
		defineWeatherTool(),
		definePopulationTool(),
		defineOrderTool(),
	}
}

func handleToolCalls(calls []openai.ToolCall, dialogue []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	for _, call := range calls {
		switch call.Function.Name {
		case "get_current_weather":
			dialogue = append(dialogue, handleWeatherCall(call))
		case "get_population":
			dialogue = append(dialogue, handlePopulationCall(call))
		case "get_user_order":
			dialogue = append(dialogue, handleOrderCall(call))
		}
	}
	return dialogue
}
