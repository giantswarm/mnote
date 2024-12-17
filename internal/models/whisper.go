package models

import (
	"fmt"

	"github.com/giantswarm/mnote/internal/config"
)

const (
	// DefaultEnglishModel is the default model for English transcription
	DefaultEnglishModel = "faster-whisper-medium-en-cpu"
	// DefaultLargeModel is the universal model for other languages
	DefaultLargeModel = "systran-faster-whisper-large-v3"
)

// WhisperModel represents a language-specific Whisper model configuration
type WhisperModel struct {
	Name     string
	Language string
	Default  bool
}

// GetWhisperModel returns the appropriate Whisper model for the given language
func GetWhisperModel(lang string, cfg *config.Config) (string, error) {
	model := cfg.GetWhisperModel(lang)
	if model == "" {
		return "", fmt.Errorf("no model found for language: %s", lang)
	}
	return model, nil
}

// ValidateLanguage checks if the given language is supported
func ValidateLanguage(lang string) bool {
	switch lang {
	case "auto", "en", "de", "es", "fr":
		return true
	default:
		return false
	}
}

// GetSupportedLanguages returns a list of supported languages
func GetSupportedLanguages() []string {
	return []string{"auto", "en", "de", "es", "fr"}
}
