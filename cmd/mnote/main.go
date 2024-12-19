package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"net/http"

	"github.com/giantswarm/mnote/internal/config"
	"github.com/giantswarm/mnote/internal/interfaces"
	"github.com/giantswarm/mnote/internal/process"
	"github.com/giantswarm/mnote/internal/registry"
	"github.com/giantswarm/mnote/internal/summarize"
	"github.com/giantswarm/mnote/internal/transcribe"
	"github.com/giantswarm/mnote/internal/utils"
	"github.com/giantswarm/mnote/internal/whisper"
	"github.com/spf13/cobra"
)

// Options holds the command-line options
type Options struct {
	VideoDir     string
	PromptName   string
	Language     string
	ForceRebuild bool
}

// usageError represents an error that should trigger usage information
type usageError struct {
	msg string
}

func (e *usageError) Error() string {
	return e.msg
}

func init() {
	// Register transcription backends
	registry.RegisterBackend("kubeai", func(cfg *config.Config) (interfaces.Transcriber, error) {
		return &transcribe.KubeAITranscriber{
			Config: cfg,
			Client: &http.Client{},
		}, nil
	})

	registry.RegisterBackend("local", func(cfg *config.Config) (interfaces.Transcriber, error) {
		return whisper.NewLocalWhisper(cfg)
	})
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
		return &usageError{fmt.Sprintf("invalid language: %s (supported: auto, en, de, es, fr)", opts.Language)}
	}

	// Validate prompt file
	promptDir := filepath.Join(os.Getenv("HOME"), ".config", "mnote", "prompts")
	promptFile := filepath.Join(promptDir, opts.PromptName)
	if _, err := os.Stat(promptFile); os.IsNotExist(err) {
		return fmt.Errorf("prompt file does not exist: %s", promptFile)
	}

	// Initialize components
	var transcriber interfaces.Transcriber
	var summarizer summarize.Summarizer

	// Initialize transcriber based on configuration
	transcriber, err = transcribe.NewTranscriber(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize transcriber: %w", err)
	}

	summarizer, err = summarize.NewSummarizer(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize summarizer: %w", err)
	}

	processor := process.NewProcessor(cfg, transcriber, summarizer)

	// Process video files in directory
	fmt.Printf("Processing videos in: %s\n", opts.VideoDir)
	fmt.Printf("Using language: %s\n", opts.Language)
	fmt.Printf("Using prompt: %s\n", opts.PromptName)
	fmt.Printf("Force rebuild: %v\n", opts.ForceRebuild)

	// Create process options
	processOpts := process.Options{
		Language:     opts.Language,
		PromptName:   opts.PromptName,
		ForceRebuild: opts.ForceRebuild,
	}

	// Process all video files in the directory
	entries, err := os.ReadDir(opts.VideoDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	foundVideo := false
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filePath := filepath.Join(opts.VideoDir, entry.Name())
		if utils.IsVideoFile(filePath) {
			foundVideo = true
			if err := processor.ProcessVideo(filePath, processOpts); err != nil {
				return fmt.Errorf("failed to process video: %w", err)
			}
		}
	}

	if !foundVideo {
		return fmt.Errorf("no supported video files found in directory: %s", opts.VideoDir)
	}

	return nil
}

// isUsageError determines if an error is related to command usage
func isUsageError(err error) bool {
	if _, ok := err.(*usageError); ok {
		return true
	}
	return strings.Contains(err.Error(), "unknown command") ||
		strings.Contains(err.Error(), "unknown flag") ||
		strings.Contains(err.Error(), "invalid argument") ||
		strings.Contains(err.Error(), "accepts")
}

func main() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		if isUsageError(err) {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			cmd.Usage()
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
