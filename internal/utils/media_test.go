package utils

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractAudio(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create a minimal valid MP4 file for testing
	videoPath := filepath.Join(tmpDir, "test.mp4")
	videoContent := []byte{
		0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, // ftyp box
		0x69, 0x73, 0x6F, 0x6D, 0x00, 0x00, 0x02, 0x00, // isom brand
		0x69, 0x73, 0x6F, 0x6D, 0x69, 0x73, 0x6F, 0x32, // compatible brands
		0x61, 0x76, 0x63, 0x31, 0x6D, 0x70, 0x34, 0x31, // compatible brands
	}
	if err := os.WriteFile(videoPath, videoContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create mock ffmpeg runner
	mock := &MockFFmpegRunner{}
	origFFmpeg := defaultFFmpeg
	defaultFFmpeg = mock
	defer func() { defaultFFmpeg = origFFmpeg }()

	// Test with supported format
	audioPath, err := ExtractAudio(videoPath, false)
	if err != nil {
		t.Errorf("ExtractAudio() error = %v", err)
	}
	if filepath.Dir(audioPath) != tmpDir {
		t.Errorf("ExtractAudio() output not in source directory")
	}
	if filepath.Ext(audioPath) != ".mp3" {
		t.Errorf("ExtractAudio() wrong output format")
	}

	// Get initial file info
	initialContent, err := os.ReadFile(audioPath)
	if err != nil {
		t.Fatal(err)
	}
	initialInfo, err := os.Stat(audioPath)
	if err != nil {
		t.Fatal(err)
	}

	// Reset mock and test with existing file and no force rebuild
	mock.ExtractCalled = false
	audioPath2, err := ExtractAudio(videoPath, false)
	if err != nil {
		t.Errorf("ExtractAudio() error = %v", err)
	}
	if audioPath2 != audioPath {
		t.Errorf("ExtractAudio() should return same path for existing file")
	}
	if mock.ExtractCalled {
		t.Error("ExtractAudio() should not call ffmpeg when file exists and forceRebuild is false")
	}

	// Verify file was not regenerated
	currentContent, err := os.ReadFile(audioPath)
	if err != nil {
		t.Fatal(err)
	}
	currentInfo, err := os.Stat(audioPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(currentContent, initialContent) {
		t.Error("ExtractAudio() should not modify file content when forceRebuild is false")
	}
	if !currentInfo.ModTime().Equal(initialInfo.ModTime()) {
		t.Error("ExtractAudio() should not modify file when forceRebuild is false")
	}

	// Test with force rebuild
	mock.ExtractCalled = false
	audioPath3, err := ExtractAudio(videoPath, true)
	if err != nil {
		t.Errorf("ExtractAudio() error = %v", err)
	}
	if audioPath3 != audioPath {
		t.Errorf("ExtractAudio() should return same path even with force rebuild")
	}
	if !mock.ExtractCalled {
		t.Error("ExtractAudio() should call ffmpeg when forceRebuild is true")
	}

	// Test with unsupported format
	unsupportedPath := filepath.Join(tmpDir, "test.xyz")
	if err := os.WriteFile(unsupportedPath, []byte("dummy content"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ExtractAudio(unsupportedPath, false); err == nil {
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
