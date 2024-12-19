package registry

import (
	"testing"

	"github.com/giantswarm/mnote/internal/config"
	"github.com/giantswarm/mnote/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

type mockTranscriber struct{}

func (m *mockTranscriber) TranscribeAudio(audioPath string, lang string) (string, error) {
	return "mock transcription", nil
}

func TestBackendRegistry(t *testing.T) {
	// Create test config
	cfg := &config.Config{
		TranscriptionAPIURL: "http://test.local",
		DefaultLanguage:    "auto",
		WhisperModelEN:    "faster-whisper-medium-en-cpu",
		WhisperModelDE:    "systran-faster-whisper-large-v3",
		WhisperModelES:    "systran-faster-whisper-large-v3",
		WhisperModelFR:    "systran-faster-whisper-large-v3",
		ChatGPTModel:      "gpt-4o",
	}

	t.Run("register and retrieve backend", func(t *testing.T) {
		// Register a test backend
		RegisterBackend("test", func(cfg *config.Config) (interfaces.Transcriber, error) {
			return &mockTranscriber{}, nil
		})

		// Retrieve the backend
		backend, err := GetBackend("test", cfg)
		assert.NoError(t, err)
		assert.NotNil(t, backend)
		assert.IsType(t, &mockTranscriber{}, backend)
	})

	t.Run("unknown backend", func(t *testing.T) {
		// Try to get a non-existent backend
		backend, err := GetBackend("nonexistent", cfg)
		assert.Error(t, err)
		assert.Nil(t, backend)
		assert.Contains(t, err.Error(), "unknown transcription backend")
	})
}
