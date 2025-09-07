package main

import (
	"fmt"
	"context"
	"os"
	"encoding/json"
	"github.com/openai/openai-go/v2"
	"github.com/anthropics/anthropic-sdk-go"

	claudeOption "github.com/anthropics/anthropic-sdk-go/option"
	openaiOption "github.com/openai/openai-go/v2/option"
)

func handleMessageProcessingAnthropic(message string) (string, error){
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	client := anthropic.NewClient(
		claudeOption.WithAPIKey(apiKey),
	)

	systemPromptBytes, err := os.ReadFile("system_prompt.txt")
	if err != nil {
		systemPromptBytes = []byte("You are a helpful assistant")
	}
	systemPrompt := string(systemPromptBytes)

	response, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5SonnetLatest,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewAssistantMessage(anthropic.NewTextBlock(systemPrompt)),
			anthropic.NewUserMessage(anthropic.NewTextBlock(message)),
		},
	})
	if err != nil {
		return "", err
	}

	text, err := extractText([]byte(response.JSON.Content.Raw()))
	if err != nil {
		return "", err
	}

	fmt.Printf("Anthropic response: '%s'\n", text)
	return text, nil
}

func HandleMessageProcessingOpenAI(message string) (string, error){
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(openaiOption.WithAPIKey(apiKey))

	fmt.Printf("Processing '%s' as a message\n", message)
	response, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(""),
			openai.UserMessage(message),
		},
	})
	if err != nil {
		return "", err
	}

	fmt.Printf("OpenAI response: '%s'\n", response.Choices[0].Message.Content)
	return string(response.Choices[0].Message.Content), nil
}

type Message struct {
    Type string `json:"type"`
    Text string `json:"text"`
}

func extractText(raw []byte) (string, error) {
    var messages []Message
    if err := json.Unmarshal(raw, &messages); err != nil {
        return "", err
    }

    if len(messages) > 0 {
        return messages[0].Text, nil
    }

    return "", nil
}
