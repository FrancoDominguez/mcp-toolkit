package main

import (
	"net/http"
	"log"
	"github.com/joho/godotenv"
	"os"
	"fmt"
)

var systemPromptFolder = os.Getenv("SYSTEM_PROMPT_FOLDER_NAME")

var app App = App{
	chatSessionId: "my_session",
	systemPromptPath: fmt.Sprintf("%s/%s", systemPromptFolder, "helpful_agent.txt"),
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

    http.HandleFunc("/ws", wsHandler)
    http.HandleFunc("/webhook", webhookHandler)
    log.Println("WebSocket server started on ws://localhost:8080/ws")
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
       log.Fatal("Error starting server:", err)
    }
}
