package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.DefaultLanguage != "auto" {
		t.Errorf("expected DefaultLanguage to be 'auto', got %s", cfg.DefaultLanguage)
	}

	expectedModels := map[string]string{
		"en": "faster-whisper-medium-en-cpu",
		"de": "systran-faster-whisper-large-v3",
		"es": "systran-faster-whisper-large-v3",
		"fr": "systran-faster-whisper-large-v3",
	}

	for lang, expectedModel := range expectedModels {
		if model, ok := cfg.WhisperModels[lang]; !ok || model != expectedModel {
			t.Errorf("expected WhisperModels[%s] to be %s, got %s", lang, expectedModel, model)
		}
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
			want:     "systran-faster-whisper-large-v3", // Updated to use large model
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
			name:     "unsupported language",
			language: "invalid",
			want:     "systran-faster-whisper-large-v3", // Updated to use large model as fallback
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
	configContent := `TRANSCRIPTION_API_URL=https://test.api/transcribe
DEFAULT_LANGUAGE=de
WHISPER_MODEL_EN=custom-en-model
CHATGPT_MODEL=gpt-4-turbo`

	err = os.WriteFile(filepath.Join(configDir, "config"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Verify loaded values
	if cfg.TranscriptionAPIURL != "https://test.api/transcribe" {
		t.Errorf("expected TranscriptionAPIURL to be 'https://test.api/transcribe', got %s", cfg.TranscriptionAPIURL)
	}
	if cfg.DefaultLanguage != "de" {
		t.Errorf("expected DefaultLanguage to be 'de', got %s", cfg.DefaultLanguage)
	}
	if cfg.WhisperModels["en"] != "custom-en-model" {
		t.Errorf("expected WhisperModels[en] to be 'custom-en-model', got %s", cfg.WhisperModels["en"])
	}
	if cfg.ChatGPTModel != "gpt-4-turbo" {
		t.Errorf("expected ChatGPTModel to be 'gpt-4-turbo', got %s", cfg.ChatGPTModel)
	}
}
