package config

import (
	"os"
	"time"
)

type Config struct {
	Shell            string
	HistoryPath      string
	PreviewCommands  bool
	DangerousPrompts bool
	EnableAI         bool
	GroqAPIKey       string
	GroqBaseURL      string
	GroqModel        string
	AIRequestTimeout time.Duration
}

func Default() Config {
	apiKey := os.Getenv("GROQ_API_KEY")
	baseURL := os.Getenv("GROQ_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.groq.com/openai/v1"
	}

	model := os.Getenv("GROQ_MODEL")
	if model == "" {
		model = "llama-3.3-70b-versatile"
	}

	return Config{
		Shell:            "cmd",
		HistoryPath:      "data/history.txt",
		PreviewCommands:  true,
		DangerousPrompts: true,
		EnableAI:         apiKey != "",
		GroqAPIKey:       apiKey,
		GroqBaseURL:      baseURL,
		GroqModel:        model,
		AIRequestTimeout: 15 * time.Second,
	}
}
