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
	case "ssp":
		fmt.Printf("Setting system prompt to: %s\n", args[0])
		err := agentConfig.SetSystemPrompt(args[0])
		if err != nil {
			return "", fmt.Errorf("ErrorHandlingCommand: %w", err)
		}
		return fmt.Sprintf("System prompt set to: %s\n", args[0]), nil
	case "gsp":
		return fmt.Sprintf("System prompt: '%s'\n", agentConfig.SystemPrompt), nil
	case "sch":
		return fmt.Sprintf("Chat history not implemented yet\n"), nil
	default:
		fmt.Printf("Unknown command: %s\n", command)
		return "", ErrUnknownCommand
	}
}
