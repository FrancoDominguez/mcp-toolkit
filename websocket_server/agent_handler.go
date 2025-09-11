package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/openai/openai-go/v2"
	claudeOption "github.com/anthropics/anthropic-sdk-go/option"
	openaiOption "github.com/openai/openai-go/v2/option"
)

type AgentResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

type Message struct {
    Type string `json:"type"`
    Text string `json:"text"`
}

func handleMessageProcessingAnthropic(message string) (string, error){
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	client := anthropic.NewClient(
		claudeOption.WithAPIKey(apiKey),
	)

	systemPrompt := fetchSystemPrompt()

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

	systemPrompt := fetchSystemPrompt()

	fmt.Printf("Processing '%s' as a message\n", message)
	response, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(message),
		},
	})
	if err != nil {
		return "", err
	}

	fmt.Printf("OpenAI response: '%s'\n", response.Choices[0].Message.Content)
	return string(response.Choices[0].Message.Content), nil
}

func HandleMessageProcessingAgent(message string) (string, error){
	url := os.Getenv("AGENT_URL")
	context_db_url := "insert context here"
	system_prompt := fetchSystemPrompt()
	request_body := fmt.Sprintf("{\"system_prompt\": \"%s\", \"user_prompt\": \"%s\", \"context_db_url\": \"%s\"}", system_prompt, message, context_db_url)
	data := []byte(request_body)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/agent/prompt", url), bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response = AgentResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Message, nil
}

func extractText(raw []byte) (string, error) {
    var messages []Message
    err := json.Unmarshal(raw, &messages)
	if err != nil {
        return "", err
    }

    if len(messages) > 0 {
        return messages[0].Text, nil
    }

    return "", nil
}

func fetchSystemPrompt() string {
	systemPromptPath := os.Getenv("SYSTEM_PROMPT_PATH")
	systemPromptBytes, err := os.ReadFile(systemPromptPath)
	if err != nil {
		systemPromptBytes = []byte("You are a helpful assistant")
	}
	systemPrompt := string(systemPromptBytes)
	return systemPrompt
}