package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()

	// Test default values
	opts := &Options{}
	cmd.Flags().VisitAll(func(f *cobra.Flag) {
		switch f.Name {
		case "prompt":
			opts.PromptName = f.DefValue
		case "language":
			opts.Language = f.DefValue
		case "force":
			val, _ := cmd.Flags().GetBool(f.Name)
			opts.ForceRebuild = val
		}
	})

	if opts.PromptName != "summarize" {
		t.Errorf("expected default prompt to be 'summarize', got %s", opts.PromptName)
	}
	if opts.ForceRebuild {
		t.Error("expected force rebuild to be false by default")
	}
}

func TestRunValidation(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	videoDir := filepath.Join(tmpDir, "videos")
	os.MkdirAll(videoDir, 0755)

	// Create temporary config directory with prompt
	configDir := filepath.Join(tmpDir, ".config", "mnote")
	promptDir := filepath.Join(configDir, "prompts")
	os.MkdirAll(promptDir, 0755)
	os.WriteFile(filepath.Join(promptDir, "summarize"), []byte("test prompt"), 0644)

	// Set HOME for config loading
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	tests := []struct {
		name    string
		opts    *Options
		wantErr bool
	}{
		{
			name: "valid options",
			opts: &Options{
				VideoDir:   videoDir,
				PromptName: "summarize",
				Language:   "en",
			},
			wantErr: false,
		},
		{
			name: "invalid directory",
			opts: &Options{
				VideoDir:   "/nonexistent",
				PromptName: "summarize",
				Language:   "en",
			},
			wantErr: true,
		},
		{
			name: "invalid language",
			opts: &Options{
				VideoDir:   videoDir,
				PromptName: "summarize",
				Language:   "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid prompt",
			opts: &Options{
				VideoDir:   videoDir,
				PromptName: "nonexistent",
				Language:   "en",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
