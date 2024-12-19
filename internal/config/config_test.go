package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Test default values
	if cfg.TranscriptionAPIURL != "https://api.kubeai.org/v1/audio/transcriptions" {
		t.Errorf("expected TranscriptionAPIURL to be 'https://api.kubeai.org/v1/audio/transcriptions', got %s", cfg.TranscriptionAPIURL)
	}
	if cfg.DefaultLanguage != "auto" {
		t.Errorf("expected DefaultLanguage to be 'auto', got %s", cfg.DefaultLanguage)
	}
	if cfg.WhisperModelEN != "faster-whisper-medium-en-cpu" {
		t.Errorf("expected WhisperModelEN to be 'faster-whisper-medium-en-cpu', got %s", cfg.WhisperModelEN)
	}
	if cfg.WhisperModelDE != "systran-faster-whisper-large-v3" {
		t.Errorf("expected WhisperModelDE to be 'systran-faster-whisper-large-v3', got %s", cfg.WhisperModelDE)
	}
	if cfg.WhisperModelES != "systran-faster-whisper-large-v3" {
		t.Errorf("expected WhisperModelES to be 'systran-faster-whisper-large-v3', got %s", cfg.WhisperModelES)
	}
	if cfg.WhisperModelFR != "systran-faster-whisper-large-v3" {
		t.Errorf("expected WhisperModelFR to be 'systran-faster-whisper-large-v3', got %s", cfg.WhisperModelFR)
	}
	if cfg.ChatGPTModel != "gpt-4o" {
		t.Errorf("expected ChatGPTModel to be 'gpt-4o', got %s", cfg.ChatGPTModel)
	}
}

func TestGetWhisperModel(t *testing.T) {
	cfg := DefaultConfig()

	tests := []struct {
		name     string
		language string
		want     string
	}{
		{
			name:     "auto detection",
			language: "auto",
			want:     "faster-whisper-medium-en-cpu",
		},
		{
			name:     "english model",
			language: "en",
			want:     "faster-whisper-medium-en-cpu",
		},
		{
			name:     "german model",
			language: "de",
			want:     "systran-faster-whisper-large-v3",
		},
		{
			name:     "spanish model",
			language: "es",
			want:     "systran-faster-whisper-large-v3",
		},
		{
			name:     "french model",
			language: "fr",
			want:     "systran-faster-whisper-large-v3",
		},
		{
			name:     "unsupported language",
			language: "invalid",
			want:     "faster-whisper-medium-en-cpu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.GetWhisperModel(tt.language)
			if got != tt.want {
				t.Errorf("GetWhisperModel(%s) = %s, want %s", tt.language, got, tt.want)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create temporary config directory
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	configDir := filepath.Join(tmpDir, ".config", "mnote")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Test with custom config
	configContent := `TRANSCRIPTION_API_URL=https://api.kubeai.org/v1/audio/transcriptions
DEFAULT_LANGUAGE=en
WHISPER_MODEL_EN=faster-whisper-medium-en-cpu
WHISPER_MODEL_DE=systran-faster-whisper-large-v3
WHISPER_MODEL_ES=systran-faster-whisper-large-v3
WHISPER_MODEL_FR=systran-faster-whisper-large-v3
CHATGPT_MODEL=gpt-4o`

	err = os.WriteFile(filepath.Join(configDir, "config"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Verify loaded values
	if cfg.TranscriptionAPIURL != "https://api.kubeai.org/v1/audio/transcriptions" {
		t.Errorf("expected TranscriptionAPIURL to be 'https://api.kubeai.org/v1/audio/transcriptions', got %s", cfg.TranscriptionAPIURL)
	}
	if cfg.DefaultLanguage != "en" {
		t.Errorf("expected DefaultLanguage to be 'en', got %s", cfg.DefaultLanguage)
	}
	if cfg.WhisperModelEN != "faster-whisper-medium-en-cpu" {
		t.Errorf("expected WhisperModelEN to be 'faster-whisper-medium-en-cpu', got %s", cfg.WhisperModelEN)
	}
	if cfg.ChatGPTModel != "gpt-4o" {
		t.Errorf("expected ChatGPTModel to be 'gpt-4o', got %s", cfg.ChatGPTModel)
	}
}
