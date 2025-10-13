package main

import (
	"fmt"
	"net/http"
	"sync"
	"errors"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
       return true
    },
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
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
			fmt.Printf("Error: %s", err)
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
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
			commandOutputMessage, err := HandleCommand(agentConfig, message[1:])
			if err != nil {
				fmt.Printf("Error: %s", err)
				switch {
				case errors.Is(err, ErrUnknownCommand):
					return "Unknown command", nil
				case errors.Is(err, ErrSystemPromptNotFound):
					return "System prompt not found", nil
				default:
					return "", fmt.Errorf("ErrProcessingCommand: %w", err)
				}
			}
			processedResponse = commandOutputMessage
		case '#':
			prompt := message[1:]
			displayMessage := fmt.Sprintf("Processing agent chat: '%s'", prompt)
			fmt.Println(displayMessage)
			conn.WriteMessage(websocket.TextMessage, []byte(displayMessage))
			agentResponse, err := handleLlmCallCustomAgent(prompt)
			if err != nil {
				return "", fmt.Errorf("ErrProcessingAgentChat: %w", err)
			}
			processedResponse = agentResponse

		default:
			fmt.Printf("Default message processing: '%s'", message)
			conn.WriteMessage(websocket.TextMessage, []byte("Default message processing: " + message))
			processedResponse = "Default message has been processed: " + message
		}
	return processedResponse, nil
}
