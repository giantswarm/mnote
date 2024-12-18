package whisper

import (
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
