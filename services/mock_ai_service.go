package services

import (
	"context"
	"eatwo/models"
	"math/rand"
)

type MockAIService struct{}

func NewMockAIService() *MockAIService {
	return &MockAIService{}
}

var mockClaudeResponse = []models.ClaudeResponse{
	{
		Content: []models.ClaudeContent{
			{Text: "Dreaming of teeth falling out is a common dream that can symbolize feelings of powerlessness or fear of aging. It may also reflect concerns about your appearance or how others perceive you.", Type: "text"},
		},
		ID:           "msg_013Zva2CMHLNnXjNJJKqJ2EF",
		Model:        "claude-3-5-sonnet-20240620",
		Role:         "assistant",
		StopReason:   "end_turn",
		StopSequence: nil,
		Type:         "message",
		Usage: models.ClaudeUsage{
			InputTokens:  2095,
			OutputTokens: 503,
		},
	},
	{
		Content: []models.ClaudeContent{
			{Text: "Dreaming of being chased can indicate that you are avoiding a situation or a person in your waking life. It may also suggest that you are feeling threatened or pressured.", Type: "text"},
		},
		ID:           "msg_014Zva2CMHLNnXjNJJKqJ2EG",
		Model:        "claude-3-5-sonnet-20240621",
		Role:         "assistant",
		StopReason:   "end_turn",
		StopSequence: nil,
		Type:         "message",
		Usage: models.ClaudeUsage{
			InputTokens:  2100,
			OutputTokens: 510,
		},
	},
	{
		Content: []models.ClaudeContent{
			{Text: "Dreaming of water can represent your emotional state. Calm and clear water might suggest peace and tranquility, while turbulent water could indicate stress or emotional turmoil.", Type: "text"},
		},
		ID:           "msg_014Zva2CMHLNnXjNJJKqJ2EG",
		Model:        "claude-3-5-sonnet-20240621",
		Role:         "assistant",
		StopReason:   "end_turn",
		StopSequence: nil,
		Type:         "message",
		Usage: models.ClaudeUsage{
			InputTokens:  2100,
			OutputTokens: 510,
		},
	},
	{
		Content: []models.ClaudeContent{
			{Text: "Dreaming of falling is often associated with feelings of insecurity or loss of control. It may reflect anxieties about failing or losing status in your waking life.", Type: "text"},
		},
		ID:           "msg_014Zva2CMHLNnXjNJJKqJ2EG",
		Model:        "claude-3-5-sonnet-20240621",
		Role:         "assistant",
		StopReason:   "end_turn",
		StopSequence: nil,
		Type:         "message",
		Usage: models.ClaudeUsage{
			InputTokens:  2100,
			OutputTokens: 510,
		},
	},
	{
		Content: []models.ClaudeContent{
			{Text: "Dreaming of flying often symbolizes a sense of freedom and liberation. It may indicate that you are feeling free from constraints or that you are seeking to escape from a situation in your waking life.", Type: "text"},
		},
		ID:           "msg_014Zva2CMHLNnXjNJJKqJ2EG",
		Model:        "claude-3-5-sonnet-20240621",
		Role:         "assistant",
		StopReason:   "end_turn",
		StopSequence: nil,
		Type:         "message",
		Usage: models.ClaudeUsage{
			InputTokens:  2100,
			OutputTokens: 510,
		},
	},
}

func (c *MockAIService) Call(ctx context.Context, prompt string) (*models.AIResponse, error) {
	i := rand.Intn(len(mockClaudeResponse))
	resp := models.AIResponse{
		Content: mockClaudeResponse[i].Content[0].Text,
	}

	return &resp, nil
}
