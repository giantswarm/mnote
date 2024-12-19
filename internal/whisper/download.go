package whisper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/giantswarm/mnote/internal/config"
)

// parseHFURL parses a HuggingFace URL in the format "hf://owner/model"
func parseHFURL(url string) (string, string, error) {
	if !strings.HasPrefix(url, "hf://") {
		return "", "", fmt.Errorf("invalid HuggingFace URL format: %s", url)
	}

	parts := strings.Split(strings.TrimPrefix(url, "hf://"), "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid HuggingFace URL format: %s", url)
	}

	return parts[0], parts[1], nil
}

// expandPath expands the ~ in the path to the user's home directory
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

// DownloadModel downloads a model from HuggingFace and saves it to the specified path
func DownloadModel(cfg *config.Config, lang string) (string, error) {
	// Get model based on language
	modelName := cfg.GetWhisperModel(lang)
	if modelName == "" {
		return "", fmt.Errorf("no model configured for language: %s", lang)
	}

	// Construct model path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	modelPath := filepath.Join(homeDir, ".config", "mnote", "models", modelName+".bin")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create model directory: %w", err)
	}

	// Check if model already exists
	if _, err := os.Stat(modelPath); err == nil {
		return modelPath, nil // Model already exists
	}

	// Construct HuggingFace download URL
	var downloadURL string
	switch modelName {
	case "faster-whisper-medium-en-cpu":
		downloadURL = "https://huggingface.co/systran/faster-whisper-medium-en/resolve/main/model.bin"
	case "systran-faster-whisper-large-v3":
		downloadURL = "https://huggingface.co/systran/faster-whisper-large-v3/resolve/main/model.bin"
	default:
		return "", fmt.Errorf("unknown model: %s", modelName)
	}

	// Download model
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to download model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download model: HTTP %d", resp.StatusCode)
	}

	// Create temporary file for downloading
	tmpFile := modelPath + ".download"
	f, err := os.Create(tmpFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile) // Clean up temporary file on error

	// Copy data to temporary file
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return "", fmt.Errorf("failed to save model: %w", err)
	}
	f.Close()

	// Move temporary file to final location
	if err := os.Rename(tmpFile, modelPath); err != nil {
		return "", fmt.Errorf("failed to move model to final location: %w", err)
	}

	return modelPath, nil
}
