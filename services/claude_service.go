package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ---
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestPayload struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system"`
}

type ClaudeService struct{}

func NewClaudeService() *ClaudeService {
	return &ClaudeService{}
}

func (c *ClaudeService) Call(ctx context.Context, prompt string) (string, error) {
	url := "https://api.anthropic.com/v1/messages"
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	payload := RequestPayload{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 1024,
		System:    "You are master of dream understanding and have access to all knowledge of dream book. Explain the meaning of a dream for a user ",
		Messages: []Message{
			{Role: "user", Content: "In my dream I was running away from a bear but could not run."},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("anthropic-beta", "message-batches-2024-09-24")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, %v", resp.StatusCode, string(body))
	}

	println(string(body))

	return string(body), nil
}
