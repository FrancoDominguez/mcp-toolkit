package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AgentResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

func handleLlmCallCustomAgent(message string) (string, error){
	url := os.Getenv("AGENT_URL")
	requestBody := map[string]string{
		"system_prompt": agentConfig.SystemPrompt,
		"user_prompt": message,
		"conversation_id": agentConfig.ConversationId,
	}

	data, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	
	fmt.Printf("Processing '%s' as a message\n", message)
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
