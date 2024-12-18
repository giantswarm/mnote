package whisper

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetModelPath(t *testing.T) {
	tests := []struct {
		name        string
		modelName   string
		wantErr    bool
		checkExists bool
	}{
		{
			name:      "valid model name",
			modelName: "test-model",
			wantErr:  false,
		},
		{
			name:      "empty model name",
			modelName: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := GetModelPath(tt.modelName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetModelPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && path == "" {
				t.Error("GetModelPath() returned empty path")
			}
		})
	}
}

func TestGetDefaultModel(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Override ModelsDir for testing
	originalModelsDir := ModelsDir
	ModelsDir = tmpDir
	defer func() { ModelsDir = originalModelsDir }()

	tests := []struct {
		name     string
		lang     string
		wantSize string
	}{
		{
			name:     "english language",
			lang:     "en",
			wantSize: "base",
		},
		{
			name:     "auto detection",
			lang:     "auto",
			wantSize: "large",
		},
		{
			name:     "other language",
			lang:     "de",
			wantSize: "large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := GetDefaultModel(tt.lang)
			if model.Size != tt.wantSize {
				t.Errorf("GetDefaultModel() size = %v, want %v", model.Size, tt.wantSize)
			}
			if model.URL == "" {
				t.Error("GetDefaultModel() returned empty URL")
			}
		})
	}
}

func TestDownloadModel(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "whisper-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override ModelsDir for testing
	originalModelsDir := ModelsDir
	ModelsDir = filepath.Join(tmpDir, "models")
	defer func() { ModelsDir = originalModelsDir }()

	testModel := ModelInfo{
		Name:     "test-model",
		Size:     "tiny",
		Language: "en",
		URL:      "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.en.bin",
	}

	if err := DownloadModel(testModel); err != nil {
		t.Errorf("DownloadModel() error = %v", err)
	}

	// Check if model file exists
	modelPath, err := GetModelPath(testModel.Name)
	if err != nil {
		t.Errorf("GetModelPath() error = %v", err)
	}

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Error("Model file was not created")
	}
}
