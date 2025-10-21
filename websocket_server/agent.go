package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrSystemPromptNotFound = errors.New("SystemPromptNotFound")
var ErrSystemPromptReadError = errors.New("ReadEror")

type Agent struct {
	Name string
	SystemPrompt string
	ConversationId string
}

func (a *Agent) SetSystemPrompt(systemPromptName string) error {
	systemPrompt, err := fetchSystemPrompt(systemPromptName)
	if err != nil {
		switch {
		case errors.Is(err, ErrSystemPromptNotFound):
			return ErrSystemPromptNotFound
		case errors.Is(err, ErrSystemPromptReadError):
			return ErrSystemPromptReadError
		default:
			return err
		}
	}

	fmt.Printf("System prompt set to: %s", systemPrompt)
	a.SystemPrompt = systemPrompt
	return nil
}

func (a *Agent) SetConversationHistory(conversation_id string) {
	a.ConversationId = conversation_id
}

func fetchSystemPrompt(systemPromptName string) (string, error) {
	systemPromptFolderPath := os.Getenv("SYSTEM_PROMPT_FOLDER_PATH")
	if systemPromptFolderPath == "" {
		systemPromptFolderPath = "./system_prompts"
	}

	systemPromptPath := fmt.Sprintf("%s/%s.txt", systemPromptFolderPath, systemPromptName)
	if _, err := os.Stat(systemPromptPath); os.IsNotExist(err) {
		return "", ErrSystemPromptNotFound
	}

	systemPromptBytes, err := os.ReadFile(systemPromptPath)
	if err != nil {
		return "", ErrSystemPromptReadError
	}
	return string(systemPromptBytes), nil
}
