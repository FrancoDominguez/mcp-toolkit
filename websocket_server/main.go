package main

import (
	"net/http"
	"fmt"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
    http.HandleFunc("/ws", wsHandler)
    http.HandleFunc("/webhook", webhookHandler)
    fmt.Println("WebSocket server started on ws://localhost:8080/ws")
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
       fmt.Println("Error starting server:", err)
    }
}
