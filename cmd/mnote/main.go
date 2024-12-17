package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/giantswarm/mnote/internal/config"
	"github.com/spf13/cobra"
)

type Options struct {
	VideoDir     string
	PromptName   string
	Language     string
	ForceRebuild bool
}

func NewRootCmd() *cobra.Command {
	opts := &Options{
		PromptName: "summarize",
		Language:   "",
	}

	cmd := &cobra.Command{
		Use:   "mnote [flags] video-directory",
		Short: "Process video files to generate transcriptions and summaries",
		Long: `mnote is a tool for processing video files to generate transcriptions and summaries.
It supports multiple languages and custom prompts for summarization.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.VideoDir = args[0]
			return run(opts)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&opts.PromptName, "prompt", "p", opts.PromptName,
		"Name of the prompt file to use for summarization")
	cmd.Flags().StringVarP(&opts.Language, "language", "l", opts.Language,
		"Language of the audio (en, de, es, fr, auto)")
	cmd.Flags().BoolVarP(&opts.ForceRebuild, "force", "f", false,
		"Force rebuild of transcription and summary")

	return cmd
}

func run(opts *Options) error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate video directory
	if _, err := os.Stat(opts.VideoDir); os.IsNotExist(err) {
		return fmt.Errorf("video directory does not exist: %s", opts.VideoDir)
	}

	// Set default language from config if not specified
	if opts.Language == "" {
		opts.Language = cfg.DefaultLanguage
	}

	// Validate language
	validLangs := map[string]bool{
		"auto": true,
		"en":   true,
		"de":   true,
		"es":   true,
		"fr":   true,
	}
	if !validLangs[opts.Language] {
		return fmt.Errorf("invalid language: %s", opts.Language)
	}

	// Validate prompt file
	promptDir := filepath.Join(os.Getenv("HOME"), ".config", "mnote", "prompts")
	promptFile := filepath.Join(promptDir, opts.PromptName)
	if _, err := os.Stat(promptFile); os.IsNotExist(err) {
		return fmt.Errorf("prompt file does not exist: %s", promptFile)
	}

	// TODO: Implement video processing logic
	fmt.Printf("Processing videos in: %s\n", opts.VideoDir)
	fmt.Printf("Using language: %s\n", opts.Language)
	fmt.Printf("Using prompt: %s\n", opts.PromptName)
	fmt.Printf("Force rebuild: %v\n", opts.ForceRebuild)

	return nil
}

func main() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
