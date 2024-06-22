package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/romiras/gpt-cli/internal/models"
)

type (
	Request struct {
		Messages    []models.Message `json:"messages"`
		Model       string           `json:"model"`
		MaxTokens   int              `json:"max_tokens"`
		Temperature float64          `json:"temperature"`
		TopP        float64          `json:"top_p"`
		Stream      bool             `json:"stream"`
	}

	HyperbolicProvider struct {
		modelConfig *models.Config
		apiKey      string
		messages    []models.Message
	}

	ChatResponse struct {
		Choices []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
	}
)

const ApiKeyName string = "HYPERBOLIC_API_KEY"
const ChatCompletionsEndpoint string = "https://api.hyperbolic.xyz/v1/chat/completions"

func NewHyperbolicProvider(cfg *models.Config) (APIProvider, error) {
	apiKey := os.Getenv(ApiKeyName)
	if apiKey == "" {
		return nil, fmt.Errorf("env var HYPERBOLIC_API_KEY is empty")
	}

	return &HyperbolicProvider{
		modelConfig: cfg,
		apiKey:      apiKey,
		messages:    make([]models.Message, 0),
	}, nil
}

func (hp *HyperbolicProvider) AddContext(message models.Message) {
	hp.messages = append(hp.messages, message)
}

func (hp *HyperbolicProvider) GetAnswer(chatQuestion string) (string, error) {
	hp.AddContext(models.Message{
		Role:    models.UserRole,
		Content: chatQuestion,
	})

	resp, err := hp.sendRequest()
	if err != nil {
		return "", err
	}

	answer, err := hp.parseResponse(resp)
	if err != nil {
		return "", err
	}

	hp.AddContext(models.Message{
		Role:    models.AssistantRole,
		Content: answer,
	})

	return answer, nil
}

func (hp *HyperbolicProvider) buildRequest() *Request {
	req := &Request{
		Model:       hp.modelConfig.Model,
		MaxTokens:   hp.modelConfig.MaxTokens,
		Temperature: hp.modelConfig.Temperature,
		TopP:        hp.modelConfig.TopP,
		Stream:      false,
	}

	req.Messages = make([]models.Message, 0, len(hp.messages))
	for _, message := range hp.messages {
		req.Messages = append(req.Messages, message)
	}

	return req
}

func (hp *HyperbolicProvider) sendRequest() ([]byte, error) {
	reqObj := hp.buildRequest()

	jsonBytes, err := json.Marshal(reqObj)
	if err != nil {
		return nil, err
	}

	// create a new HTTP request
	req, err := http.NewRequest("POST", ChatCompletionsEndpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+hp.apiKey)

	// create a new HTTP client
	client := &http.Client{}

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (hp *HyperbolicProvider) parseResponse(resp []byte) (string, error) {
	var chatResponse ChatResponse

	err := json.Unmarshal(resp, &chatResponse)
	if err != nil {
		return "", err
	}

	answer := ""
	for _, choice := range chatResponse.Choices {
		answer += choice.Message.Content

		if choice.FinishReason == "stop" {
			break
		}
	}

	return answer, nil
}
