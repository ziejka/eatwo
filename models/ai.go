package models

type AIResponse struct {
	Content string `json:"prompt"`
}

type ClaudeResponse struct {
	Content      []ClaudeContent `json:"content"`
	ID           string          `json:"id"`
	Model        string          `json:"model"`
	Role         string          `json:"role"`
	StopReason   string          `json:"stop_reason"`
	StopSequence *string         `json:"stop_sequence"`
	Type         string          `json:"type"`
	Usage        ClaudeUsage     `json:"usage"`
}

type ClaudeContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
