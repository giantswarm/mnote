package transcribe

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/giantswarm/mnote/internal/config"
	"github.com/giantswarm/mnote/internal/interfaces"
	"github.com/giantswarm/mnote/internal/registry"
)

// MockTranscriber implements the Transcriber interface for testing
type MockTranscriber struct {
	ReturnText string
	ReturnErr  error
}

func (m *MockTranscriber) TranscribeAudio(audioPath, language string) (string, error) {
	return m.ReturnText, m.ReturnErr
}

func init() {
	// Register mock backend for testing
	registry.RegisterBackend("kubeai", func(cfg *config.Config) (interfaces.Transcriber, error) {
		return &MockTranscriber{
			ReturnText: "Test transcription",
			ReturnErr:  nil,
		}, nil
	})

	// Register local backend with error handling
	registry.RegisterBackend("local", func(cfg *config.Config) (interfaces.Transcriber, error) {
		return &MockTranscriber{
			ReturnText: "Test transcription",
			ReturnErr:  nil,
		}, nil
	})
}

func TestNewTranscriber(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &config.Config{
				TranscriptionAPIURL: "http://example.com",
				DefaultLanguage:    "auto",
				WhisperModelEN:    "faster-whisper-medium-en-cpu",
				WhisperModelDE:    "systran-faster-whisper-large-v3",
				WhisperModelES:    "systran-faster-whisper-large-v3",
				WhisperModelFR:    "systran-faster-whisper-large-v3",
				ChatGPTModel:      "gpt-4o",
			},
			wantErr: false,
		},
		{
			name: "missing API URL",
			cfg: &config.Config{
				DefaultLanguage: "auto",
				WhisperModelEN: "faster-whisper-medium-en-cpu",
				ChatGPTModel:   "gpt-4o",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transcriber, err := NewTranscriber(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTranscriber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && transcriber == nil {
				t.Error("NewTranscriber() returned nil transcriber without error")
			}
		})
	}
}

func TestKubeAITranscribeAudio(t *testing.T) {
	// Create test config
	cfg := &config.Config{
		TranscriptionAPIURL: "http://example.com",
		DefaultLanguage:     "auto",
		WhisperModelEN:     "faster-whisper-medium-en-cpu",
		WhisperModelDE:     "systran-faster-whisper-large-v3",
		WhisperModelES:     "systran-faster-whisper-large-v3",
		WhisperModelFR:     "systran-faster-whisper-large-v3",
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Verify content type is multipart
		contentType := r.Header.Get("Content-Type")
		if contentType == "" || contentType[:len("multipart/form-data")] != "multipart/form-data" {
			t.Errorf("expected multipart/form-data content type, got %s", contentType)
		}

		// Parse multipart form
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("failed to parse multipart form: %v", err)
		}

		// Verify required fields
		if _, ok := r.MultipartForm.File["file"]; !ok {
			t.Error("file field missing from request")
		}

		// Verify model field matches expected model for language
		model := r.FormValue("model")
		if model == "" {
			t.Error("model field missing from request")
		}
		language := r.FormValue("language")
		expectedModel := cfg.GetWhisperModel(language)
		if model != expectedModel {
			t.Errorf("incorrect model for language %s, got %s, want %s", language, model, expectedModel)
		}

		// Return test response
		result := TranscriptionResult{Text: "Test transcription"}
		json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()

	// Update config with test server URL
	cfg.TranscriptionAPIURL = server.URL

	// Create transcriber
	transcriber, err := NewTranscriber(cfg)
	if err != nil {
		t.Fatalf("failed to create transcriber: %v", err)
	}

	// Create temporary audio file
	tmpDir := t.TempDir()
	audioPath := filepath.Join(tmpDir, "test.mp3")
	if err := os.WriteFile(audioPath, []byte("test audio data"), 0644); err != nil {
		t.Fatalf("failed to create test audio file: %v", err)
	}

	tests := []struct {
		name     string
		language string
		wantErr  bool
	}{
		{
			name:     "auto detection",
			language: "auto",
			wantErr:  false,
		},
		{
			name:     "english transcription",
			language: "en",
			wantErr:  false,
		},
		{
			name:     "german transcription",
			language: "de",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := transcriber.TranscribeAudio(audioPath, tt.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranscribeAudio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && text != "Test transcription" {
				t.Errorf("TranscribeAudio() = %v, want %v", text, "Test transcription")
			}
		})
	}
}
