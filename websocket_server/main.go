package main

import (
	"net/http"
	"log"
	"github.com/joho/godotenv"
)

var agentConfig *AgentConfig

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	agentConfig = NewAgentConfig()

    http.HandleFunc("/ws", websocketHandler)
    log.Println("WebSocket server started on ws://localhost:8080/ws")
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
       log.Fatal("Error starting server:", err)
    }
}
