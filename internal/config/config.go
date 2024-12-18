package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration settings for mnote
type Config struct {
	TranscriptionAPIURL  string            `mapstructure:"TRANSCRIPTION_API_URL"`
	TranscriptionBackend string            `mapstructure:"TRANSCRIPTION_BACKEND"` // "local" or "kubeai"
	LocalModelPath       string            `mapstructure:"LOCAL_MODEL_PATH"`
	LocalModelSize       string            `mapstructure:"LOCAL_MODEL_SIZE"` // tiny, base, small, medium, large
	DefaultLanguage      string            `mapstructure:"DEFAULT_LANGUAGE"`
	WhisperModels        map[string]string `mapstructure:"-"`
	ChatGPTModel         string            `mapstructure:"CHATGPT_MODEL"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		TranscriptionAPIURL:  "https://api.kubeai.org/v1/audio/transcriptions",
		TranscriptionBackend: "kubeai", // Default to KubeAI for backward compatibility
		LocalModelPath:       "~/.config/mnote/models/ggml-base.en.bin",
		LocalModelSize:       "base",
		DefaultLanguage:      "auto",
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

	// Check for test environment
	if os.Getenv("HOME") == "/home/ubuntu/mnote/test" {
		v.SetConfigType("yaml")
		v.AddConfigPath("/home/ubuntu/mnote/test")
	} else {
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
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Create default config file only in non-test environment
			if os.Getenv("HOME") != "/home/ubuntu/mnote/test" {
				defaultConfig := `TRANSCRIPTION_API_URL=https://api.kubeai.org/v1/audio/transcriptions
TRANSCRIPTION_BACKEND=kubeai
LOCAL_MODEL_PATH=~/.config/mnote/models/ggml-base.en.bin
LOCAL_MODEL_SIZE=base
DEFAULT_LANGUAGE=auto
WHISPER_MODEL_EN=faster-whisper-medium-en-cpu
WHISPER_MODEL_DE=systran-faster-whisper-large-v3
WHISPER_MODEL_ES=systran-faster-whisper-large-v3
WHISPER_MODEL_FR=systran-faster-whisper-large-v3
CHATGPT_MODEL=gpt-4o`
				configFile := filepath.Join(os.Getenv("HOME"), ".config", "mnote", "config")
				if err := os.WriteFile(configFile, []byte(defaultConfig), 0644); err != nil {
					return nil, fmt.Errorf("failed to create config file: %w", err)
				}
				// Reload config after creating file
				if err := v.ReadInConfig(); err != nil {
					return nil, fmt.Errorf("failed to read new config file: %w", err)
				}
			} else {
				return nil, fmt.Errorf("failed to read test config file: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand model path
	if strings.HasPrefix(config.LocalModelPath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		config.LocalModelPath = filepath.Join(home, config.LocalModelPath[2:])
	}

	// Validate configuration
	if err := ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
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
		// For auto-detection, use the large model
		return c.WhisperModels["de"]
	}
	if model, ok := c.WhisperModels[lang]; ok {
		return model
	}
	// Fallback to large model if language not supported
	return c.WhisperModels["de"]
}
