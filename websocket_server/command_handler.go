package main
import (
	"strings"
	"fmt"
)

func HandleCommand(message string){
	commandArgs := strings.Split(message, "/")
	prefix := commandArgs[0]
	switch prefix {
	case "chat-history":
		command := commandArgs[1]
		switch command {
			case "set-history":
				fmt.Printf("Setting chat history")
		}
	}
}