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
		TranscriptionAPIURL: "https://api.kubeai.org/v1/audio/transcriptions",
		DefaultLanguage: "auto",
		WhisperModels: map[string]string{
			"en": "faster-whisper-medium-en-cpu",
			"de": "systran-faster-whisper-large-v3",
			"es": "systran-faster-whisper-large-v3",
			"fr": "systran-faster-whisper-large-v3",
		},
		ChatGPTModel: "gpt-4o",
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
	promptsDir := filepath.Join(configDir, "prompts")
	v.AddConfigPath(configDir)

	// Create config directory and prompts directory
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directories: %w", err)
	}

	// Create default prompt file if it doesn't exist
	defaultPromptFile := filepath.Join(promptsDir, "summarize")
	if _, err := os.Stat(defaultPromptFile); os.IsNotExist(err) {
		defaultPrompt := `Create a detailed summary of the following meeting transcript. Structure the summary according to the main topics discussed and organize the information into logical sections. For each topic, summarize who was involved, what was discussed in detail, what decisions were made, what problems or challenges were identified, and what solutions were proposed or implemented.`
		if err := os.WriteFile(defaultPromptFile, []byte(defaultPrompt), 0644); err != nil {
			return nil, fmt.Errorf("failed to create default prompt file: %w", err)
		}
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Create default config file
			defaultConfig := `TRANSCRIPTION_API_URL=https://api.kubeai.org/v1/audio/transcriptions
DEFAULT_LANGUAGE=auto
WHISPER_MODEL_EN=faster-whisper-medium-en-cpu
WHISPER_MODEL_DE=systran-faster-whisper-large-v3
WHISPER_MODEL_ES=systran-faster-whisper-large-v3
WHISPER_MODEL_FR=systran-faster-whisper-large-v3
CHATGPT_MODEL=gpt-4o`
			configFile := filepath.Join(configDir, "config")
			if err := os.WriteFile(configFile, []byte(defaultConfig), 0644); err != nil {
				return nil, fmt.Errorf("failed to create config file: %w", err)
			}
			// Reload config after creating file
			if err := v.ReadInConfig(); err != nil {
				return nil, fmt.Errorf("failed to read new config file: %w", err)
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
