package main

import (
	"fmt"
	"context"
	"os"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"

)

func HandleMessageProcessing(message []byte) (string, error){
	apiKey := os.Getenv("OPENAI_API_KEY")
	fmt.Println("API Key: ", apiKey)
	client := openai.NewClient(option.WithAPIKey(apiKey))
	messageString := string(message)

	fmt.Printf("Processing '%s' as a message\n", messageString)
	chat_completion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful assistant."),
			openai.UserMessage(messageString),
		},
	})
	if err != nil {
		return "woops", err
	}

	fmt.Printf("OpenAI response: '%s'\n", chat_completion.Choices[0].Message.Content)
	return string(chat_completion.Choices[0].Message.Content), nil
}