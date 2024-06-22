package api

import "github.com/romiras/gpt-cli/internal/models"

type (
	APIProvider interface {
		GetAnswer(chatQuestion string) (string, error)
		AddContext(message models.Message)
	}
)
