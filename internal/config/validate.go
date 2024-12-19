package config

import (
	"fmt"
)

// ValidateConfig validates the configuration settings
func ValidateConfig(cfg *Config) error {
	// Validate required fields
	if cfg.TranscriptionAPIURL == "" {
		return fmt.Errorf("TRANSCRIPTION_API_URL is required")
	}

	// Validate language models
	if cfg.WhisperModelEN == "" {
		return fmt.Errorf("WHISPER_MODEL_EN is required")
	}
	if cfg.WhisperModelDE == "" {
		return fmt.Errorf("WHISPER_MODEL_DE is required")
	}
	if cfg.WhisperModelES == "" {
		return fmt.Errorf("WHISPER_MODEL_ES is required")
	}
	if cfg.WhisperModelFR == "" {
		return fmt.Errorf("WHISPER_MODEL_FR is required")
	}

	// Validate default language
	validLanguages := map[string]bool{
		"auto": true,
		"en":   true,
		"de":   true,
		"es":   true,
		"fr":   true,
	}
	if !validLanguages[cfg.DefaultLanguage] {
		return fmt.Errorf("invalid DEFAULT_LANGUAGE: %s (must be auto, en, de, es, or fr)", cfg.DefaultLanguage)
	}

	// Validate ChatGPT model
	if cfg.ChatGPTModel == "" {
		return fmt.Errorf("CHATGPT_MODEL is required")
	}

	return nil
}
