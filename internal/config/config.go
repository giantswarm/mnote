package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration settings for mnote
type Config struct {
	TranscriptionAPIURL string            `mapstructure:"TRANSCRIPTION_API_URL"`
	DefaultLanguage    string            `mapstructure:"DEFAULT_LANGUAGE"`
	WhisperModels     map[string]string `mapstructure:"-"`
	ChatGPTModel      string            `mapstructure:"CHATGPT_MODEL"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		DefaultLanguage: "auto",
		WhisperModels: map[string]string{
			"en": "faster-whisper-medium-en-cpu",
			"de": "systran-faster-whisper-large-v3",
			"es": "systran-faster-whisper-large-v3",
			"fr": "systran-faster-whisper-large-v3",
		},
		ChatGPTModel: "gpt-4",
	}
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Set up viper
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("env")

	// Config file path
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "mnote")
	v.AddConfigPath(configDir)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Create default config file if it doesn't exist
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create config directory: %w", err)
			}
			configFile := filepath.Join(configDir, "config")
			if err := os.WriteFile(configFile, []byte(""), 0644); err != nil {
				return nil, fmt.Errorf("failed to create config file: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Load language-specific Whisper models from environment variables
	languages := []string{"en", "de", "es", "fr"}
	for _, lang := range languages {
		envKey := fmt.Sprintf("WHISPER_MODEL_%s", lang)
		if model := v.GetString(envKey); model != "" {
			config.WhisperModels[lang] = model
		}
	}

	return config, nil
}

// GetWhisperModel returns the appropriate Whisper model for the given language
func (c *Config) GetWhisperModel(lang string) string {
	if lang == "auto" {
		// For auto-detection, use the English model as default
		return c.WhisperModels["en"]
	}
	if model, ok := c.WhisperModels[lang]; ok {
		return model
	}
	// Fallback to English model if language not supported
	return c.WhisperModels["en"]
}
