package utils

import (
	"fmt"
	"path/filepath"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// SupportedVideoFormats contains the list of supported video file extensions
var SupportedVideoFormats = []string{".mp4", ".mkv", ".avi", ".mov"}

// ExtractAudio extracts audio from a video file and saves it in the same directory
func ExtractAudio(videoPath string) (string, error) {
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

	// Extract audio using ffmpeg
	err := ffmpeg.Input(videoPath).
		Output(audioPath, ffmpeg.KwArgs{
			"acodec": "libmp3lame",
			"ab":     "192k",
			"ar":     "44100",
			"y":      "", // Overwrite output file if it exists
		}).
		OverWriteOutput().
		Run()

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
