package main

import (
	"encoding/json"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type Item struct {
	Name     string
	Quantity int
	Price    float64
}

type Order struct {
	UserID    string
	Items     []Item
	Total     float64
	CreatedAt time.Time
}

func defineOrderTool() openai.Tool {
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"userID": {
				Type:        jsonschema.String,
				Description: "用戶ID",
			},
		},
		Required: []string{"userID"},
	}
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "get_user_order",
			Description: "獲取指定用戶的訂單",
			Parameters:  params,
		},
	}
}

func getUserOrder(userID string) Order {
	return Order{
		UserID: userID,
		Items: []Item{
			{Name: "筆記本", Quantity: 2, Price: 15.99},
			{Name: "鉛筆", Quantity: 10, Price: 0.99},
		},
		Total:     47.88,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}
}

func handleOrderCall(call openai.ToolCall) openai.ChatCompletionMessage {
	var params struct {
		UserID string `json:"userID"`
	}
	json.Unmarshal([]byte(call.Function.Arguments), &params)
	order := getUserOrder(params.UserID)
	orderJSON, _ := json.Marshal(order)
	return openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    string(orderJSON),
		Name:       call.Function.Name,
		ToolCallID: call.ID,
	}
}
