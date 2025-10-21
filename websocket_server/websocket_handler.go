package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
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

		log.Println("\n\n========================================")
		response, err := processMessage(message)
		if err != nil {
			log.Println("Error: ", err)
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
		}
		log.Println("Response: ", response)
		conn.WriteMessage(websocket.TextMessage, []byte(response))
	}
}

func processMessage(message string) (string, error){
	log.Printf("Received: %s\n", message)
	firstChar := message[0]
	switch firstChar {
		case '/':
			log.Println("Processing command")
			commandOutputMessage, err := HandleCommand(message[1:])
			if err != nil {
				log.Println("Error: ", err)
				return "", err
			}
			return commandOutputMessage, nil
		case '#':
			log.Println("Processing agent chat")
			prompt := message[1:]
			agentResponse, err := handleLlmCallCustomAgent(prompt)
			if err != nil {
				return "", err
			}
			return agentResponse, nil

		default:
			return "Message received", nil
		}
}