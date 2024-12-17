package process

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/giantswarm/mnote/internal/config"
	"github.com/giantswarm/mnote/internal/summarize"
	"github.com/giantswarm/mnote/internal/transcribe"
)

type mockTranscriber struct {
	transcript string
	err       error
}

func (m *mockTranscriber) TranscribeAudio(audioPath, model, language string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.transcript, nil
}

type mockSummarizer struct {
	summary string
	err    error
}

func (m *mockSummarizer) SummarizeTranscript(transcript, promptName string) (string, error) {
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

	// Create mock dependencies
	cfg := config.DefaultConfig()
	transcriber := &mockTranscriber{transcript: "Test transcript"}
	summarizer := &mockSummarizer{summary: "Test summary"}

	// Create processor
	processor := NewProcessor(cfg, transcriber, summarizer)

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
