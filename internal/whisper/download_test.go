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
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	// Create test config
	cfg := &config.Config{
		DefaultLanguage: "auto",
		WhisperModelEN: "faster-whisper-medium-en-cpu",
		WhisperModelDE: "systran-faster-whisper-large-v3",
		WhisperModelES: "systran-faster-whisper-large-v3",
		WhisperModelFR: "systran-faster-whisper-large-v3",
	}

	tests := []struct {
		name        string
		lang        string
		wantModel   string
		setupMock   func(t *testing.T, modelPath string)
		wantErr     bool
		errContains string
	}{
		{
			name:      "english model",
			lang:      "en",
			wantModel: "faster-whisper-medium-en-cpu",
			setupMock: func(t *testing.T, modelPath string) {
				if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
					t.Fatalf("Failed to create model directory: %v", err)
				}
				if err := os.WriteFile(modelPath, []byte("test model data"), 0644); err != nil {
					t.Fatalf("Failed to write test model: %v", err)
				}
			},
			wantErr: false,
		},
		{
			name:      "german model",
			lang:      "de",
			wantModel: "systran-faster-whisper-large-v3",
			setupMock: func(t *testing.T, modelPath string) {
				if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
					t.Fatalf("Failed to create model directory: %v", err)
				}
				if err := os.WriteFile(modelPath, []byte("test model data"), 0644); err != nil {
					t.Fatalf("Failed to write test model: %v", err)
				}
			},
			wantErr: false,
		},
		{
			name:        "unsupported language",
			lang:        "invalid",
			wantModel:   "faster-whisper-medium-en-cpu",
			wantErr:     false,
			setupMock: func(t *testing.T, modelPath string) {
				if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
					t.Fatalf("Failed to create model directory: %v", err)
				}
				if err := os.WriteFile(modelPath, []byte("test model data"), 0644); err != nil {
					t.Fatalf("Failed to write test model: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelPath := filepath.Join(tmpDir, ".config", "mnote", "models", tt.wantModel+".bin")
			if tt.setupMock != nil {
				tt.setupMock(t, modelPath)
			}

			path, err := DownloadModel(cfg, tt.lang)
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
			if path != modelPath {
				t.Errorf("path = %q, want %q", path, modelPath)
			}
		})
	}
}
