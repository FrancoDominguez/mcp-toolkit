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
		_, messageSlice, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		message := string(messageSlice)

		response, err := processMessage(conn, message)
		if err != nil {
			fmt.Println("Error processing message:", err)
		}
		fmt.Println("Response: ", response)
		conn.WriteMessage(websocket.TextMessage, []byte(response))
	}
}

func processMessage(conn *websocket.Conn, message string) (string, error){

	firstChar := message[0]
	processedResponse := ""
	switch firstChar {
		case '/':
			fmt.Println("Processing command: ", message)
			conn.WriteMessage(websocket.TextMessage, []byte("Processing command: " + message))
			processedResponse = "Command has been processed: " + message
		case '#':
			fmt.Printf("Processing agent request: '%s'\n", message)
			conn.WriteMessage(websocket.TextMessage, []byte("Processing agent request: " + message))
			response, err := HandleMessageProcessingAgent(message)
			if err != nil {
				fmt.Println("Error processing message:", err)
			}
			processedResponse = response
		default:
			fmt.Println("Default message processing: ", message)
			conn.WriteMessage(websocket.TextMessage, []byte("Default message processing: " + message))
			processedResponse = "Default message has been processed: " + message
		}
	return processedResponse, nil
}