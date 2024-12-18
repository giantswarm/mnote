package transcribe

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/giantswarm/mnote/internal/config"
)

func TestTranscribeAudio(t *testing.T) {
	// Create test config
	cfg := config.DefaultConfig()

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
		if model := r.FormValue("model"); model == "" {
			t.Error("model field missing from request")
		}

		// Return test response
		result := TranscriptionResult{Text: "Test transcription"}
		json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()

	// Update config with test server URL
	cfg.TranscriptionAPIURL = server.URL

	// Create transcriber
	transcriber := NewTranscriber(cfg)

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
			result, err := transcriber.TranscribeAudio(audioPath, tt.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranscribeAudio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Text != "Test transcription" {
				t.Errorf("TranscribeAudio() = %v, want %v", result.Text, "Test transcription")
			}
		})
	}
}
