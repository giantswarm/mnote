package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractAudio(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create a dummy video file
	videoPath := filepath.Join(tmpDir, "test.mp4")
	if err := os.WriteFile(videoPath, []byte("dummy video content"), 0644); err != nil {
		t.Fatalf("Failed to create test video file: %v", err)
	}

	// Test with supported format
	audioPath, err := ExtractAudio(videoPath)
	if err != nil {
		t.Errorf("ExtractAudio() error = %v", err)
	}
	if filepath.Dir(audioPath) != tmpDir {
		t.Errorf("ExtractAudio() output not in source directory")
	}
	if filepath.Ext(audioPath) != ".mp3" {
		t.Errorf("ExtractAudio() wrong output format")
	}

	// Test with unsupported format
	unsupportedPath := filepath.Join(tmpDir, "test.xyz")
	if err := os.WriteFile(unsupportedPath, []byte("dummy content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if _, err := ExtractAudio(unsupportedPath); err == nil {
		t.Error("ExtractAudio() should fail with unsupported format")
	}
}

func TestIsVideoFile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"mp4 file", "video.mp4", true},
		{"mkv file", "video.mkv", true},
		{"avi file", "video.avi", true},
		{"mov file", "video.mov", true},
		{"unsupported format", "video.xyz", false},
		{"no extension", "video", false},
		{"case insensitive", "video.MP4", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVideoFile(tt.path); got != tt.want {
				t.Errorf("IsVideoFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
