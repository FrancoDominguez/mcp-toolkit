package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrSystemPromptNotFound = errors.New("SystemPromptNotFound")
var ErrSystemPromptReadError = errors.New("ReadEror")

type AgentConfig struct {
	Name string
	SystemPrompt string
	ConversationId string
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		Name: "Jarvis",
		SystemPrompt: "default",
		ConversationId: NewUUID(),
	}
}

func (a *AgentConfig) SetSystemPrompt(systemPromptName string) error {
	systemPrompt, err := fetchSystemPrompt(systemPromptName)
	if err != nil {
		return err
	}
	a.SystemPrompt = systemPrompt
	return nil
}

func (a *AgentConfig) SetConversationHistory(conversationId string) {
	a.ConversationId = conversationId
}

func (a *AgentConfig) NewConversation() (string) {
	conversationId := NewUUID()
	a.SetConversationHistory(conversationId)
	return conversationId
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