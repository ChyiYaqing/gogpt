package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseURL              string          `yaml:"base_url"`
	Model                string          `yaml:"model"`
	Timeout              time.Duration   `yaml:"timeout"`
	MaxConversationTurns int             `yaml:"max_conversation_turns"`
	StreamBufferSize     int             `yaml:"stream_buffer_size"`
	KeepFullResponse     bool            `yaml:"keep_full_response"`
	ModelParams          ModelParameters `yaml:"model_params"`
	LogLevel             string          `yaml:"log_level"`
	LogFile              string          `yaml:"log_file"`
}

type ModelParameters struct {
	Temperature   float64 `yaml:"temperature"`
	TopP          float64 `yaml:"top_p"`
	TopK          int     `yaml:"top_k"`
	RepeatPenalty float64 `yaml:"repeat_penalty"`
	MaxTokens     int     `yaml:"max_tokens"`
}

func Load(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return LoadDefault(), nil
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadDefault() *Config {
	return &Config{
		BaseURL:              "http://localhost:11434",
		Model:                "llama3.2",
		Timeout:              30 * time.Second,
		MaxConversationTurns: 10,
		StreamBufferSize:     4096,
		KeepFullResponse:     false,
		ModelParams: ModelParameters{
			Temperature:   0.7,
			TopP:          0.9,
			TopK:          40,
			RepeatPenalty: 1.1,
			MaxTokens:     2000,
		},
		LogLevel: "info",
		LogFile:  "ollama-chat.log",
	}
}
