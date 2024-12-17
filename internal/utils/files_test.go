package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetOutputPath(t *testing.T) {
	tests := []struct {
		name      string
		inputPath string
		suffix    string
		want      string
	}{
		{
			name:      "with suffix",
			inputPath: "/path/to/video.mp4",
			suffix:    "summary",
			want:      "/path/to/video_summary.md",
		},
		{
			name:      "without suffix",
			inputPath: "/path/to/video.mp4",
			suffix:    "",
			want:      "/path/to/video.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOutputPath(tt.inputPath, tt.suffix); got != tt.want {
				t.Errorf("GetOutputPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"existing file", testFile, true},
		{"non-existing file", filepath.Join(tmpDir, "nonexistent.txt"), false},
		{"directory", tmpDir, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.path); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "subdir", "test.txt")
	testData := []byte("test data")

	if err := WriteFile(testPath, testData); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Verify file exists and contains correct data
	got, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}
	if string(got) != string(testData) {
		t.Errorf("WriteFile() wrote %v, want %v", string(got), string(testData))
	}
}
