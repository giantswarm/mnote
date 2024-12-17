package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// FFmpegRunner defines the interface for audio extraction
type FFmpegRunner interface {
	ExtractAudioFromVideo(inputPath, outputPath string) error
}

// DefaultFFmpegRunner implements FFmpegRunner using ffmpeg-go
type DefaultFFmpegRunner struct{}

// ExtractAudioFromVideo implements FFmpegRunner interface
func (r *DefaultFFmpegRunner) ExtractAudioFromVideo(inputPath, outputPath string) error {
	return ffmpeg.Input(inputPath).
		Output(outputPath, ffmpeg.KwArgs{
			"acodec": "libmp3lame",
			"ab":     "192k",
			"ar":     "44100",
			"y":      "", // Overwrite output file if it exists
		}).
		OverWriteOutput().
		Run()
}

// MockFFmpegRunner implements FFmpegRunner for testing
type MockFFmpegRunner struct {
	ExtractCalled bool
	ForceError    bool
}

func (m *MockFFmpegRunner) ExtractAudioFromVideo(inputPath, outputPath string) error {
	m.ExtractCalled = true
	if m.ForceError {
		return fmt.Errorf("mock ffmpeg error")
	}
	return os.WriteFile(outputPath, []byte("mock mp3 content"), 0644)
}

// defaultFFmpeg is the default FFmpeg runner implementation
var defaultFFmpeg FFmpegRunner = &DefaultFFmpegRunner{}

// SetFFmpegRunner allows setting a custom FFmpeg runner (used for testing)
func SetFFmpegRunner(runner FFmpegRunner) {
	defaultFFmpeg = runner
}

// SupportedVideoFormats contains the list of supported video file extensions
var SupportedVideoFormats = []string{".mp4", ".mkv", ".avi", ".mov"}

// ExtractAudio extracts audio from a video file and saves it in the same directory
func ExtractAudio(videoPath string, forceRebuild bool) (string, error) {
	// Validate video format
	ext := strings.ToLower(filepath.Ext(videoPath))
	supported := false
	for _, format := range SupportedVideoFormats {
		if ext == format {
			supported = true
			break
		}
	}
	if !supported {
		return "", fmt.Errorf("unsupported video format: %s", ext)
	}

	// Generate output path in the same directory
	dir := filepath.Dir(videoPath)
	base := strings.TrimSuffix(filepath.Base(videoPath), ext)
	audioPath := filepath.Join(dir, base+".mp3")

	// Check if file exists and skip if not forcing rebuild
	if !forceRebuild && FileExists(audioPath) {
		fmt.Printf("Audio file already exists: %s\n", audioPath)
		return audioPath, nil
	}

	// Extract audio using ffmpeg
	err := defaultFFmpeg.ExtractAudioFromVideo(videoPath, audioPath)
	if err != nil {
		return "", fmt.Errorf("failed to extract audio: %w", err)
	}

	return audioPath, nil
}

// IsVideoFile checks if the given file is a supported video file
func IsVideoFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, format := range SupportedVideoFormats {
		if ext == format {
			return true
		}
	}
	return false
}
