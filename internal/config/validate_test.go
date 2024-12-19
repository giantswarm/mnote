package config

import "testing"

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				TranscriptionAPIURL: "https://api.example.com",
				DefaultLanguage:    "auto",
				WhisperModelEN:    "faster-whisper-medium-en-cpu",
				WhisperModelDE:    "systran-faster-whisper-large-v3",
				WhisperModelES:    "systran-faster-whisper-large-v3",
				WhisperModelFR:    "systran-faster-whisper-large-v3",
				ChatGPTModel:      "gpt-4o",
			},
			wantErr: false,
		},
		{
			name: "missing transcription URL",
			cfg: &Config{
				DefaultLanguage: "auto",
				WhisperModelEN: "faster-whisper-medium-en-cpu",
				ChatGPTModel:   "gpt-4o",
			},
			wantErr: true,
		},
		{
			name: "invalid language",
			cfg: &Config{
				TranscriptionAPIURL: "https://api.example.com",
				DefaultLanguage:    "invalid",
				WhisperModelEN:    "faster-whisper-medium-en-cpu",
				ChatGPTModel:      "gpt-4o",
			},
			wantErr: true,
		},
		{
			name: "missing ChatGPT model",
			cfg: &Config{
				TranscriptionAPIURL: "https://api.example.com",
				DefaultLanguage:    "auto",
				WhisperModelEN:    "faster-whisper-medium-en-cpu",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
