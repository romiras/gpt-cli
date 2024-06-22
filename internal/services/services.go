package services

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/romiras/gpt-cli/internal/api"
	"github.com/romiras/gpt-cli/internal/models"
)

func Run(provider api.APIProvider) error {
	messages := make([]models.Message, 0)
	for {
		userInput := getUserInput()

		if userInput == "" {
			break
		}

		messages = append(messages, models.Message{
			Role:    models.UserRole,
			Content: userInput,
		})
		answer, err := provider.GetAnswer(userInput)

		if err != nil {
			fmt.Println(err.Error())
			break
		}

		fmt.Printf("Answer: %s\n\n", answer)
	}

	return nil
}

func getUserInput() string {
	fmt.Print("Enter your question (multiple lines allowed, Ctrl+D to finish): ")
	reader := bufio.NewReader(os.Stdin)
	var userInput string

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		userInput += string(line) + "\n"
	}

	return strings.TrimSuffix(userInput, "\n")
}
