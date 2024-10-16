package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {
	godotenv.Load("./.env")
	ctx := context.Background()
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	tools := defineTools()

	dialogue := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleUser, Content: `{"userID": "ray", "msg": 能幫我查我的訂單嗎？}`},
	}

	fmt.Println("初始對話：")
	printMessages(dialogue)

	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: dialogue,
			Tools:    tools,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion 錯誤: %v\n", err)
		return
	}

	fmt.Println("\nAI 的初始回應：")
	printMessage(resp.Choices[0].Message)

	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		fmt.Println("\n檢測到工具調用：")
		for _, call := range resp.Choices[0].Message.ToolCalls {
			fmt.Printf("工具名稱: %s\n參數: %s\n", call.Function.Name, call.Function.Arguments)
		}

		dialogue = append(dialogue, resp.Choices[0].Message)

		dialogue = handleToolCalls(resp.Choices[0].Message.ToolCalls, dialogue)

		fmt.Println("\n更新後的對話歷史：")
		printMessages(dialogue)

		resp, err = client.CreateChatCompletion(ctx,
			openai.ChatCompletionRequest{
				Model:    openai.GPT4TurboPreview,
				Messages: dialogue,
			},
		)

		if err != nil {
			fmt.Printf("第二次 ChatCompletion 錯誤: %v\n", err)
			return
		}
	}

	fmt.Println("\nAI 的最終回答：")
	printMessage(resp.Choices[0].Message)
}

func printMessages(messages []openai.ChatCompletionMessage) {
	for i, msg := range messages {
		fmt.Printf("消息 %d:\n", i+1)
		printMessage(msg)
		fmt.Println()
	}
}

func printMessage(msg openai.ChatCompletionMessage) {
	fmt.Printf("角色: %s\n內容: %s\n", msg.Role, msg.Content)
	if len(msg.ToolCalls) > 0 {
		fmt.Println("工具調用:")
		for _, call := range msg.ToolCalls {
			fmt.Printf("  名稱: %s\n  參數: %s\n", call.Function.Name, call.Function.Arguments)
		}
	}

}
