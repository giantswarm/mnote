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
				TranscriptionAPIURL: "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name: "valid local config",
			cfg: &Config{
				TranscriptionBackend: "local",
				LocalModelPath:      "~/.config/mnote/models/model.bin",
				LocalModelSize:      "base",
			},
			wantErr: false,
		},
		{
			name: "invalid backend",
			cfg: &Config{
				TranscriptionBackend: "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid model size",
			cfg: &Config{
				TranscriptionBackend: "local",
				LocalModelPath:      "path/to/model",
				LocalModelSize:      "invalid",
			},
			wantErr: true,
		},
		{
			name: "missing model path",
			cfg: &Config{
				TranscriptionBackend: "local",
				LocalModelSize:      "base",
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
