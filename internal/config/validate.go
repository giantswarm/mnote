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

	// Validate catalog configuration
	if len(cfg.Catalog) == 0 {
		return fmt.Errorf("catalog configuration is required")
	}

	for name, model := range cfg.Catalog {
		// Validate required fields
		if model.URL == "" {
			return fmt.Errorf("model %s: URL is required", name)
		}
		if model.Owner == "" {
			return fmt.Errorf("model %s: owner is required", name)
		}
		if model.Engine == "" {
			return fmt.Errorf("model %s: engine is required", name)
		}
		if len(model.Features) == 0 {
			return fmt.Errorf("model %s: at least one feature is required", name)
		}

		// Validate local path for local backend
		if cfg.TranscriptionBackend == "local" && model.LocalPath == "" {
			return fmt.Errorf("model %s: local path is required when using local transcription backend", name)
		}

		// Validate engine type
		if model.Engine != "FasterWhisper" {
			return fmt.Errorf("model %s: invalid engine %s (must be 'FasterWhisper')", name, model.Engine)
		}

		// Validate features
		validFeatures := map[string]bool{"SpeechToText": true}
		for _, feature := range model.Features {
			if !validFeatures[feature] {
				return fmt.Errorf("model %s: invalid feature %s (must be 'SpeechToText')", name, feature)
			}
		}
	}

	return nil
}
