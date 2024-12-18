package process

import (
	"fmt"

	"github.com/giantswarm/mnote/internal/config"
	"github.com/giantswarm/mnote/internal/interfaces"
	"github.com/giantswarm/mnote/internal/summarize"
	"github.com/giantswarm/mnote/internal/utils"
)

// Options holds the processing options
type Options struct {
	Language     string
	PromptName   string
	ForceRebuild bool
}

// Processor handles the complete video processing workflow
type Processor struct {
	config      *config.Config
	transcriber interfaces.Transcriber
	summarizer  summarize.Summarizer
}

// NewProcessor creates a new Processor instance
func NewProcessor(cfg *config.Config, transcriber interfaces.Transcriber, summarizer summarize.Summarizer) *Processor {
	return &Processor{
		config:      cfg,
		transcriber: transcriber,
		summarizer:  summarizer,
	}
}

// ProcessVideo processes a video file, generating transcription and summary
func (p *Processor) ProcessVideo(path string, opts Options) error {
	// Validate video file
	if !utils.IsVideoFile(path) {
		return fmt.Errorf("not a supported video file: %s", path)
	}

	// Extract audio
	audioPath, err := utils.ExtractAudio(path, opts.ForceRebuild)
	if err != nil {
		return fmt.Errorf("failed to extract audio: %w", err)
	}

	// Get output paths
	transcriptPath := utils.GetOutputPath(path, "transcript")
	summaryPath := utils.GetOutputPath(path, opts.PromptName)

	// Skip transcription if file exists and not forcing rebuild
	if !opts.ForceRebuild && utils.FileExists(transcriptPath) {
		fmt.Printf("Transcript file already exists: %s\n", transcriptPath)
	} else {
		// Perform transcription
		result, err := p.transcriber.TranscribeAudio(audioPath, opts.Language)
		if err != nil {
			return fmt.Errorf("transcription failed: %w", err)
		}

		// Save transcript
		if err := utils.WriteFile(transcriptPath, []byte(result)); err != nil {
			return fmt.Errorf("failed to save transcript: %w", err)
		}
		fmt.Printf("Transcript saved to: %s\n", transcriptPath)
	}

	// Skip summarization if file exists and not forcing rebuild
	if !opts.ForceRebuild && utils.FileExists(summaryPath) {
		fmt.Printf("Summary file already exists: %s\n", summaryPath)
		return nil
	}

	// Read transcript for summarization
	transcript, err := utils.ReadFile(transcriptPath)
	if err != nil {
		return fmt.Errorf("failed to read transcript: %w", err)
	}

	// Generate summary
	summary, err := p.summarizer.SummarizeTranscript(string(transcript), opts.PromptName, opts.ForceRebuild)
	if err != nil {
		return fmt.Errorf("summarization failed: %w", err)
	}

	// Save summary
	if err := utils.WriteFile(summaryPath, []byte(summary)); err != nil {
		return fmt.Errorf("failed to save summary: %w", err)
	}
	fmt.Printf("Summary saved to: %s\n", summaryPath)

	return nil
}
