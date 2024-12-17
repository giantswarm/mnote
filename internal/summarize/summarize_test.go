package summarize

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/giantswarm/mnote/internal/config"
)

func TestNewSummarizer(t *testing.T) {
	// Save current API key and restore it after test
	originalKey := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", originalKey)

	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid API key",
			apiKey:  "test-key",
			wantErr: false,
		},
		{
			name:    "missing API key",
			apiKey:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("OPENAI_API_KEY", tt.apiKey)
			cfg := config.DefaultConfig()
			_, err := NewSummarizer(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSummarizer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSummarizeTranscript(t *testing.T) {
	// Skip if no API key is set
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("Skipping test: OPENAI_API_KEY not set")
	}

	// Create test config directory and prompt
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "mnote", "prompts")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Create test prompt file
	promptContent := "You are a helpful assistant. Please summarize the following text."
	promptFile := filepath.Join(configDir, "test_prompt")
	if err := os.WriteFile(promptFile, []byte(promptContent), 0644); err != nil {
		t.Fatalf("failed to create prompt file: %v", err)
	}

	// Set HOME for config loading
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Create summarizer
	cfg := config.DefaultConfig()
	summarizer, err := NewSummarizer(cfg)
	if err != nil {
		t.Fatalf("failed to create summarizer: %v", err)
	}

	// Test transcript
	transcript := "This is a test transcript that needs to be summarized."

	// Test summarization
	summary, err := summarizer.SummarizeTranscript(transcript, "test_prompt", false)
	if err != nil {
		t.Fatalf("SummarizeTranscript() error = %v", err)
	}
	if summary == "" {
		t.Error("SummarizeTranscript() returned empty summary")
	}
}
