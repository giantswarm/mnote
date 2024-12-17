package summarize

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/giantswarm/mnote/internal/config"
	"github.com/sashabaranov/go-openai"
)

// Summarizer handles transcript summarization using OpenAI API
type Summarizer struct {
	client OpenAIClient
	config *config.Config
}

// OpenAIClient interface for mocking in tests
type OpenAIClient interface {
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// NewSummarizer creates a new Summarizer instance
func NewSummarizer(cfg *config.Config) (*Summarizer, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(apiKey)
	return &Summarizer{
		client: client,
		config: cfg,
	}, nil
}

// SummarizeTranscript generates a summary of the transcript using the specified prompt
func (s *Summarizer) SummarizeTranscript(transcript, promptName string, forceRebuild bool) (string, error) {
	// Read prompt file
	promptDir := filepath.Join(os.Getenv("HOME"), ".config", "mnote", "prompts")
	promptFile := filepath.Join(promptDir, promptName)
	promptContent, err := os.ReadFile(promptFile)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file: %w", err)
	}

	// Create chat completion request
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: s.config.ChatGPTModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: string(promptContent),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: transcript,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from API")
	}

	return resp.Choices[0].Message.Content, nil
}

// MockOpenAIClient implements OpenAIClient for testing
type MockOpenAIClient struct{}

func (m *MockOpenAIClient) CreateChatCompletion(_ context.Context, _ openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{
				Message: openai.ChatCompletionMessage{
					Content: "Mock summary: This is a test transcript summary.",
				},
			},
		},
	}, nil
}