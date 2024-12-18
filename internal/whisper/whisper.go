package whisper

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	whisperlib "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/giantswarm/mnote/internal/config"
	"github.com/giantswarm/mnote/internal/interfaces"
)

// LocalWhisper implements the Transcriber interface using local whisper.cpp
type LocalWhisper struct {
	modelPath string
	model     whisperlib.Model
	context   whisperlib.Context
	config    config.ModelConfig
}

// NewLocalWhisper creates a new LocalWhisper instance with the specified model configuration
func NewLocalWhisper(cfg config.ModelConfig) (*LocalWhisper, error) {
	// Download or get existing model
	modelPath, err := DownloadModel(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}

	// Load the whisper model
	model, err := whisperlib.New(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load whisper model: %w", err)
	}

	// Create context for transcription
	context, err := model.NewContext()
	if err != nil {
		model.Close()
		return nil, fmt.Errorf("failed to create whisper context: %w", err)
	}

	return &LocalWhisper{
		modelPath: modelPath,
		model:     model,
		context:   context,
		config:    cfg,
	}, nil
}

// New is deprecated, use NewLocalWhisper instead
func New(modelPath string) (*LocalWhisper, error) {
	return NewLocalWhisper(config.ModelConfig{
		LocalPath: modelPath,
	})
}

// Close releases resources associated with the model
func (w *LocalWhisper) Close() error {
	if w.model != nil {
		w.model.Close()
	}
	return nil
}

// TranscribeAudio implements the Transcriber interface
func (w *LocalWhisper) TranscribeAudio(audioPath string, lang string) (string, error) {
	// Create temporary file for PCM data
	tmpFile, err := os.CreateTemp("", "whisper-*.pcm")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Convert audio to raw PCM using ffmpeg-go
	err = ffmpeg.Input(audioPath).
		Output(tmpFile.Name(),
			ffmpeg.KwArgs{
				"acodec": "pcm_f32le",
				"ac":     1,
				"ar":     16000,
				"f":      "f32le",
			}).
		OverWriteOutput().
		Run()
	if err != nil {
		return "", fmt.Errorf("failed to convert audio: %w", err)
	}

	// Read PCM data
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read PCM data: %w", err)
	}

	// Convert bytes to float32 samples
	samples := make([]float32, len(data)/4)
	for i := 0; i < len(data); i += 4 {
		samples[i/4] = math.Float32frombits(uint32(data[i]) | uint32(data[i+1])<<8 | uint32(data[i+2])<<16 | uint32(data[i+3])<<24)
	}

	// Set language if specified
	if lang != "auto" {
		if err := w.context.SetLanguage(lang); err != nil {
			return "", fmt.Errorf("failed to set language %s: %w", lang, err)
		}
	}

	// Configure transcription settings
	w.context.SetTranslate(false)
	w.context.SetThreads(4)
	w.context.SetTokenTimestamps(true)

	// Process audio with callbacks to collect segments
	var transcription strings.Builder
	segmentCallback := func(segment whisperlib.Segment) {
		transcription.WriteString(segment.Text)
		transcription.WriteString(" ")
	}

	if err := w.context.Process(samples, segmentCallback, nil); err != nil {
		return "", fmt.Errorf("failed to process audio: %w", err)
	}

	return strings.TrimSpace(transcription.String()), nil
}

// Ensure LocalWhisper implements Transcriber interface
var _ interfaces.Transcriber = (*LocalWhisper)(nil)
