package process

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/giantswarm/mnote/internal/config"
	"github.com/giantswarm/mnote/internal/transcribe"
	"github.com/giantswarm/mnote/internal/utils"
)

// mockTranscriber implements transcribe.Transcriber interface
type mockTranscriber struct {
	transcript string
	err       error
}

func (m *mockTranscriber) TranscribeAudio(audioPath, language string) (*transcribe.TranscriptionResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &transcribe.TranscriptionResult{Text: m.transcript}, nil
}

// mockSummarizer implements summarize.Summarizer interface
type mockSummarizer struct {
	summary string
	err    error
}

func (m *mockSummarizer) SummarizeTranscript(transcript, promptName string, forceRebuild bool) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.summary, nil
}

func TestProcessVideo(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create test video file
	videoPath := filepath.Join(tmpDir, "test.mp4")
	if err := os.WriteFile(videoPath, []byte("dummy video content"), 0644); err != nil {
		t.Fatalf("Failed to create test video file: %v", err)
	}

	// Create mock FFmpeg runner
	mockFFmpeg := &utils.MockFFmpegRunner{}
	utils.SetFFmpegRunner(mockFFmpeg)
	defer utils.SetFFmpegRunner(&utils.DefaultFFmpegRunner{})

	// Create mock dependencies
	cfg := config.DefaultConfig()
	transcriber := &mockTranscriber{transcript: "Test transcript"}
	summarizer := &mockSummarizer{summary: "Test summary"}

	// Create processor
	processor := &Processor{
		config:      cfg,
		transcriber: transcriber,
		summarizer:  summarizer,
	}

	// Test processing
	opts := Options{
		Language:     "en",
		PromptName:   "test",
		ForceRebuild: true,
	}

	err := processor.ProcessVideo(videoPath, opts)
	if err != nil {
		t.Errorf("ProcessVideo() error = %v", err)
	}

	// Verify FFmpeg was called
	if !mockFFmpeg.ExtractCalled {
		t.Error("FFmpeg extraction was not called")
	}

	// Verify output files
	transcriptPath := filepath.Join(tmpDir, "test_transcript.md")
	summaryPath := filepath.Join(tmpDir, "test_test.md")

	if !fileExists(transcriptPath) {
		t.Error("Transcript file not created")
	}
	if !fileExists(summaryPath) {
		t.Error("Summary file not created")
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
