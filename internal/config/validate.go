package config

import (
	"fmt"
	"strings"
)

// ValidateConfig validates the configuration settings
func ValidateConfig(cfg *Config) error {
	// Validate TranscriptionBackend
	if cfg.TranscriptionBackend != "local" && cfg.TranscriptionBackend != "kubeai" {
		return fmt.Errorf("invalid transcription backend: %s (must be 'local' or 'kubeai')", cfg.TranscriptionBackend)
	}

	// Validate LocalModelSize
	validSizes := map[string]bool{
		"tiny":   true,
		"base":   true,
		"small":  true,
		"medium": true,
		"large":  true,
	}
	if cfg.TranscriptionBackend == "local" {
		if !validSizes[cfg.LocalModelSize] {
			return fmt.Errorf("invalid local model size: %s (must be tiny, base, small, medium, or large)", cfg.LocalModelSize)
		}
	}

	// Validate LocalModelPath for local backend
	if cfg.TranscriptionBackend == "local" && strings.TrimSpace(cfg.LocalModelPath) == "" {
		return fmt.Errorf("local model path is required when using local transcription backend")
	}

	return nil
}
