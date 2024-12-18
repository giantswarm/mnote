package config

import "testing"

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid kubeai config",
			cfg: &Config{
				TranscriptionBackend: "kubeai",
				TranscriptionAPIURL:  "https://api.example.com",
				Catalog: map[string]ModelConfig{
					"faster-whisper-medium-en-cpu": {
						Enabled:  true,
						Features: []string{"SpeechToText"},
						Owner:    "systran",
						URL:      "hf://systran/faster-whisper-medium-en",
						Engine:   "FasterWhisper",
						Language: "en",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid local config",
			cfg: &Config{
				TranscriptionBackend: "local",
				LocalModelPath:       "~/.config/mnote/models/model.bin",
				LocalModelSize:       "base",
				Catalog: map[string]ModelConfig{
					"faster-whisper-large-v3": {
						Enabled:   true,
						Features:  []string{"SpeechToText"},
						Owner:     "systran",
						URL:       "hf://systran/faster-whisper-large-v3",
						Engine:    "FasterWhisper",
						LocalPath: "~/.config/mnote/models/model.bin",
						Language:  "auto",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid backend",
			cfg: &Config{
				TranscriptionBackend: "invalid",
				Catalog: map[string]ModelConfig{},
			},
			wantErr: true,
		},
		{
			name: "invalid model size",
			cfg: &Config{
				TranscriptionBackend: "local",
				LocalModelPath:       "path/to/model",
				LocalModelSize:       "invalid",
				Catalog:             map[string]ModelConfig{},
			},
			wantErr: true,
		},
		{
			name: "missing model path",
			cfg: &Config{
				TranscriptionBackend: "local",
				LocalModelSize:       "base",
				Catalog:             map[string]ModelConfig{},
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
