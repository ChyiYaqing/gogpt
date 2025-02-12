package config

import "time"

type Config struct {
	BaseURL              string
	Model                string
	Timeout              time.Duration
	MaxConversationTurns int
	ModelParams          ModelParameters
}

type ModelParameters struct {
	Temperature   float64
	TopP          float64
	TopK          int
	RepeatPenalty float64
	MaxTokens     int
}


func LoadConfig() *Config {
	return &Config{
		BaseURL: "http://localhost:11434",
		Model: "llama3.2",
		Timeout: 30 * time.Second,
		MaxConversationTurns: 10,
		ModelParams: ModelParameters{
			Temperature: 0.7,
			TopP: 0.9,
			TopK: 40,
			RepeatPenalty: 1.1,
			MaxTokens: 2000,
		},
	}
}