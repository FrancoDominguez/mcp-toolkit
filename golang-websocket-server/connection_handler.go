package main

import (
	"sync"
	"fmt"
	"github.com/gorilla/websocket"
)

func handleConnectionStream(conn *websocket.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		chat_completion, err := HandleMessageProcessing(message)
		if err != nil {
			fmt.Println("Error processing message:", err)
		}
		conn.WriteMessage(websocket.TextMessage, []byte(chat_completion))
	}
}