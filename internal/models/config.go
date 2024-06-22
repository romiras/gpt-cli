package models

const DefaultModel string = "mistralai/Mixtral-8x7B-Instruct-v0.1"

type Config struct {
	Model       string
	MaxTokens   int
	Temperature float64
	TopP        float64
}

func BuildDefaultConfig() Config {
	return Config{
		Model:       DefaultModel,
		MaxTokens:   512,
		Temperature: 0.7,
		TopP:        0.9,
	}
}
