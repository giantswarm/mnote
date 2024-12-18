package whisper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ModelsDir is the directory where whisper models are stored
var ModelsDir = filepath.Join(os.Getenv("HOME"), ".config", "mnote", "models")

// ModelInfo contains information about a whisper model
type ModelInfo struct {
	Name     string
	Size     string // tiny, base, small, medium, large
	Language string // en, de, fr, es, auto
	URL      string
}

// GetModelPath returns the path to a model file
func GetModelPath(name string) (string, error) {
	// Expand home directory in ModelsDir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	modelsDir := strings.Replace(ModelsDir, "~", homeDir, 1)

	// Create models directory if it doesn't exist
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create models directory: %w", err)
	}

	return filepath.Join(modelsDir, fmt.Sprintf("%s.bin", name)), nil
}

// GetDefaultModel returns the default model info based on language
func GetDefaultModel(lang string) ModelInfo {
	switch lang {
	case "en":
		return ModelInfo{
			Name:     "ggml-base.en",
			Size:     "base",
			Language: "en",
			URL:      "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.en.bin",
		}
	case "auto":
		return ModelInfo{
			Name:     "ggml-large",
			Size:     "large",
			Language: "auto",
			URL:      "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-large.bin",
		}
	default:
		// For other languages, use the large model which supports multiple languages
		return ModelInfo{
			Name:     "ggml-large",
			Size:     "large",
			Language: "auto",
			URL:      "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-large.bin",
		}
	}
}
