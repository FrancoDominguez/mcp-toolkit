package main
import (
	"strings"
	"fmt"
	"errors"
)

var ErrUnknownCommand = errors.New("unknown command")

func HandleCommand(message string) (string, error){
	commandArgs := strings.Split(message, " ")
	command := commandArgs[0]
	args := commandArgs[1:]

	switch command {
	case "/ssp":
		err := agentConfig.SetSystemPrompt(args[0])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("System prompt set to: %s\n", args[0]), nil

	case "/gsp":
		return fmt.Sprintf("System prompt: '%s'\n", agentConfig.SystemPrompt), nil

	case "/sch":
		agentConfig.SetConversationHistory(args[0])
		return fmt.Sprintf("Conversation history set to: %s\n", args[0]), nil

	default:
		return "", ErrUnknownCommand
	}
}