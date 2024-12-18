package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/giantswarm/mnote/internal/utils"
)

// Set test environment
func init() {
	os.Setenv("TEST_ENV", "true")
}

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()

	// Test default values
	opts := &Options{}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		switch f.Name {
		case "prompt":
			opts.PromptName = f.DefValue
		case "language":
			opts.Language = f.DefValue
		case "force":
			val, _ := cmd.Flags().GetBool(f.Name)
			opts.ForceRebuild = val
		}
	})

	if opts.PromptName != "summarize" {
		t.Errorf("expected default prompt to be 'summarize', got %s", opts.PromptName)
	}
	if opts.ForceRebuild {
		t.Error("expected force rebuild to be false by default")
	}
}

func TestRunValidation(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	videoDir := filepath.Join(tmpDir, "videos")
	os.MkdirAll(videoDir, 0755)

	// Create test video file
	videoPath := filepath.Join(videoDir, "test.mp4")
	if err := os.WriteFile(videoPath, []byte("test video content"), 0644); err != nil {
		t.Fatalf("Failed to create test video file: %v", err)
	}

	// Set up mock FFmpeg runner
	mockFFmpeg := &utils.MockFFmpegRunner{}
	utils.SetFFmpegRunner(mockFFmpeg)
	defer utils.SetFFmpegRunner(&utils.DefaultFFmpegRunner{})

	// Set up mock HTTP server for transcription API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate language parameter
		language := r.FormValue("language")
		validLangs := map[string]bool{
			"auto": true, "en": true, "de": true, "es": true, "fr": true,
		}
		if !validLangs[language] {
			http.Error(w, "invalid language: "+language, http.StatusUnprocessableEntity)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"text": "mock transcription",
		})
	}))
	defer mockServer.Close()
	os.Setenv("TRANSCRIPTION_API_URL", mockServer.URL)
	defer os.Unsetenv("TRANSCRIPTION_API_URL")

	// Create temporary config directory with prompt
	configDir := filepath.Join(tmpDir, ".config", "mnote")
	promptDir := filepath.Join(configDir, "prompts")
	os.MkdirAll(promptDir, 0755)
	os.WriteFile(filepath.Join(promptDir, "summarize"), []byte("test prompt"), 0644)

	// Set HOME for config loading
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Set mock OpenAI API key
	os.Setenv("OPENAI_API_KEY", "test-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	// Create test config file
	configFile := filepath.Join(configDir, "config")
	os.WriteFile(configFile, []byte(`default_language: "en"
chatgpt_model: "gpt-4"
`), 0644)

	tests := []struct {
		name       string
		opts       *Options
		wantErr    bool
		wantUsage  bool
		setupFiles bool
	}{
		{
			name: "valid options",
			opts: &Options{
				VideoDir:   videoDir,
				PromptName: "summarize",
				Language:   "en",
			},
			wantErr:    false,
			wantUsage:  false,
			setupFiles: true,
		},
		{
			name: "invalid directory",
			opts: &Options{
				VideoDir:   "/nonexistent",
				PromptName: "summarize",
				Language:   "en",
			},
			wantErr:    true,
			wantUsage:  false,
			setupFiles: false,
		},
		{
			name: "invalid language",
			opts: &Options{
				VideoDir:   videoDir,
				PromptName: "summarize",
				Language:   "invalid",
			},
			wantErr:    true,
			wantUsage:  true,
			setupFiles: false,
		},
		{
			name: "invalid prompt",
			opts: &Options{
				VideoDir:   videoDir,
				PromptName: "nonexistent",
				Language:   "en",
			},
			wantErr:    true,
			wantUsage:  false,
			setupFiles: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock FFmpeg for each test
			mockFFmpeg.ExtractCalled = false
			mockFFmpeg.ForceError = false

			// Create test video file if needed
			if tt.setupFiles {
				if err := os.WriteFile(videoPath, []byte("test video content"), 0644); err != nil {
					t.Fatalf("Failed to create test video file: %v", err)
				}
			}

			err := run(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check if usage info is included in error message when expected
			if err != nil && tt.wantUsage {
				if !strings.Contains(err.Error(), "supported: auto, en, de, es, fr") {
					t.Errorf("run() error message does not contain supported languages list")
				}
			}

			// Check for duplicate error messages
			if err != nil {
				errCount := strings.Count(err.Error(), err.Error())
				if errCount > 1 {
					t.Errorf("run() error message appears %d times, want 1", errCount)
				}
			}
		})
	}
}
