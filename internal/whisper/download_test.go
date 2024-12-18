package whisper

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/giantswarm/mnote/internal/config"
)

func TestParseHFURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantOwner   string
		wantModel   string
		wantErr     bool
		errContains string
	}{
		{
			name:      "valid URL",
			url:       "hf://systran/faster-whisper-medium-en",
			wantOwner: "systran",
			wantModel: "faster-whisper-medium-en",
			wantErr:   false,
		},
		{
			name:        "invalid prefix",
			url:         "http://systran/faster-whisper-medium-en",
			wantErr:     true,
			errContains: "invalid HuggingFace URL format",
		},
		{
			name:        "missing model",
			url:         "hf://systran",
			wantErr:     true,
			errContains: "invalid HuggingFace URL format",
		},
		{
			name:        "empty URL",
			url:         "",
			wantErr:     true,
			errContains: "invalid HuggingFace URL format",
		},
		{
			name:        "too many parts",
			url:         "hf://systran/model/extra",
			wantErr:     true,
			errContains: "invalid HuggingFace URL format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, model, err := parseHFURL(tt.url)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if owner != tt.wantOwner {
				t.Errorf("owner = %q, want %q", owner, tt.wantOwner)
			}
			if model != tt.wantModel {
				t.Errorf("model = %q, want %q", model, tt.wantModel)
			}
		})
	}
}

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "path with tilde",
			path:    "~/models/whisper.bin",
			want:    filepath.Join(home, "models/whisper.bin"),
			wantErr: false,
		},
		{
			name:    "absolute path",
			path:    "/tmp/models/whisper.bin",
			want:    "/tmp/models/whisper.bin",
			wantErr: false,
		},
		{
			name:    "relative path",
			path:    "models/whisper.bin",
			want:    "models/whisper.bin",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandPath(tt.path)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expandPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownloadModel(t *testing.T) {
	tmpDir := t.TempDir()

	// Copy test model to testdata directory
	testModelPath := filepath.Join("testdata", "test-model.bin")
	testModelData, err := os.ReadFile(testModelPath)
	if err != nil {
		t.Fatalf("failed to read test model: %v", err)
	}

	tests := []struct {
		name        string
		config      config.ModelConfig
		setupMock   func(t *testing.T, cfg config.ModelConfig) string
		wantErr     bool
		errContains string
	}{
		{
			name: "missing URL",
			config: config.ModelConfig{
				LocalPath: filepath.Join(tmpDir, "model.bin"),
			},
			wantErr:     true,
			errContains: "model URL not specified",
		},
		{
			name: "invalid URL format",
			config: config.ModelConfig{
				URL:       "invalid://url",
				LocalPath: filepath.Join(tmpDir, "model.bin"),
			},
			wantErr:     true,
			errContains: "invalid HuggingFace URL format",
		},
		{
			name: "missing local path",
			config: config.ModelConfig{
				URL: "hf://systran/faster-whisper-medium-en",
			},
			wantErr:     true,
			errContains: "local path not specified",
		},
		{
			name: "valid config with test model",
			config: config.ModelConfig{
				URL:       "hf://systran/faster-whisper-medium-en",
				LocalPath: filepath.Join(tmpDir, "model.bin"),
			},
			setupMock: func(t *testing.T, cfg config.ModelConfig) string {
				// Create a minimal valid whisper model file for testing
				modelData := []byte{
					// Whisper magic number
					0x77, 0x68, 0x69, 0x73, // "whis"
					0x70, 0x65, 0x72, 0x00, // "per\0"
					// Version 1
					0x00, 0x00, 0x00, 0x01,
					// Model type (base)
					0x00, 0x00, 0x00, 0x01,
					// Add some dummy model data
					0x00, 0x01, 0x02, 0x03,
					0x04, 0x05, 0x06, 0x07,
				}
				if err := os.WriteFile(cfg.LocalPath, modelData, 0644); err != nil {
					t.Fatalf("Failed to write test model file: %v", err)
				}
				return cfg.LocalPath
			},
			wantErr: false,
		},
		{
			name: "language-specific model",
			config: config.ModelConfig{
				URL:       "hf://systran/faster-whisper-large-v3",
				LocalPath: filepath.Join(tmpDir, "large-v3.bin"),
				Language:  "de",
				Enabled:   true,
				Features:  []string{"SpeechToText"},
				Owner:     "systran",
				Engine:    "FasterWhisper",
			},
			setupMock: func(t *testing.T, cfg config.ModelConfig) string {
				if err := os.WriteFile(cfg.LocalPath, testModelData, 0644); err != nil {
					t.Fatalf("failed to write test model: %v", err)
				}
				return cfg.LocalPath
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expectedPath string
			if tt.setupMock != nil {
				expectedPath = tt.setupMock(t, tt.config)
			}

			path, err := DownloadModel(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("model file does not exist at %s", path)
			}
			if tt.setupMock != nil && path != expectedPath {
				t.Errorf("path = %q, want %q", path, expectedPath)
			}

			// Skip model initialization and transcription tests in test environment
			// These require actual model files and will be tested in integration tests
			if testing.Short() {
				return
			}
		})
	}
}
