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

	// Test catalog configuration
	expectedModels := map[string]struct {
		name     string
		language string
	}{
		"faster-whisper-medium-en-cpu": {
			name:     "faster-whisper-medium-en",
			language: "en",
		},
		"systran-faster-whisper-large-v3": {
			name:     "faster-whisper-large-v3",
			language: "auto",
		},
	}

	for modelID, expected := range expectedModels {
		model, exists := cfg.Catalog[modelID]
		if !exists {
			t.Errorf("expected model %s in catalog", modelID)
			continue
		}
		if !model.Enabled {
			t.Errorf("expected model %s to be enabled", modelID)
		}
		if model.URL != "hf://systran/"+expected.name {
			t.Errorf("expected URL 'hf://systran/%s', got '%s'", expected.name, model.URL)
		}
		if model.Engine != "FasterWhisper" {
			t.Errorf("expected Engine 'FasterWhisper', got '%s'", model.Engine)
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
			want:     "systran-faster-whisper-large-v3",
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
			want:     "systran-faster-whisper-large-v3",
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
	configContent := `TRANSCRIPTION_BACKEND=local
LOCAL_MODEL_SIZE=base
TRANSCRIPTION_API_URL=https://api.kubeai.org/v1/audio/transcriptions
DEFAULT_LANGUAGE=en
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
	if cfg.TranscriptionBackend != "local" {
		t.Errorf("expected TranscriptionBackend to be 'local', got %s", cfg.TranscriptionBackend)
	}

	model, exists := cfg.Catalog["faster-whisper-medium-en-cpu"]
	if !exists {
		t.Fatal("expected model 'faster-whisper-medium-en-cpu' in catalog")
	}

	if !model.Enabled {
		t.Error("expected model to be enabled")
	}

	if model.URL != "hf://systran/faster-whisper-medium-en" {
		t.Errorf("expected URL 'hf://systran/faster-whisper-medium-en', got '%s'", model.URL)
	}
}
