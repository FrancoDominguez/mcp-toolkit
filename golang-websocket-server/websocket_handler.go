package main

import (
	"fmt"
	"net/http"
	"sync"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
       return true
    },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var connectionWg sync.WaitGroup

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
       fmt.Println("Error upgrading:", err)
       return
    }
    defer conn.Close()

	connectionWg.Add(1)
	go handleConnectionStream(conn, &connectionWg)
	connectionWg.Wait()
}

func handleConnectionStream(conn *websocket.Conn, connectionWg *sync.WaitGroup) {
	defer connectionWg.Done()
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