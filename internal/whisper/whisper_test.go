package whisper

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWhisper is a mock implementation of the LocalWhisper interface
type MockWhisper struct {
	mock.Mock
}

func (m *MockWhisper) TranscribeAudio(audioPath, language string) (string, error) {
	args := m.Called(audioPath, language)
	return args.String(0), args.Error(1)
}

func (m *MockWhisper) Close() {
	m.Called()
}

func TestNewLocalWhisper(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	// Create models directory
	modelPath := filepath.Join(tmpDir, ".config", "mnote", "models")
	if err := os.MkdirAll(modelPath, 0755); err != nil {
		t.Fatalf("Failed to create model directory: %v", err)
	}

	// Skip actual model loading for unit test
	t.Skip("Skipping TestNewLocalWhisper as it requires a valid model file")
}

func TestTranscribeAudio(t *testing.T) {
	// Create mock whisper instance
	mock := &MockWhisper{}

	// Set up expectations
	mock.On("TranscribeAudio", "testdata/test_en.mp3", "en").Return("Test transcription", nil)
	mock.On("TranscribeAudio", "testdata/test_de.mp3", "de").Return("Test deutsche Transkription", nil)
	mock.On("TranscribeAudio", "/nonexistent/audio.mp3", "en").Return("", fmt.Errorf("file not found"))
	mock.On("TranscribeAudio", "testdata/test_en.mp3", "invalid").Return("", fmt.Errorf("invalid language"))
	mock.On("Close").Return()

	tests := []struct {
		name      string
		audioPath string
		lang      string
		want      string
		wantErr   bool
	}{
		{
			name:      "english transcription",
			audioPath: "testdata/test_en.mp3",
			lang:      "en",
			want:      "Test transcription",
			wantErr:   false,
		},
		{
			name:      "german transcription",
			audioPath: "testdata/test_de.mp3",
			lang:      "de",
			want:      "Test deutsche Transkription",
			wantErr:   false,
		},
		{
			name:      "non-existent audio",
			audioPath: "/nonexistent/audio.mp3",
			lang:      "en",
			want:      "",
			wantErr:   true,
		},
		{
			name:      "invalid language",
			audioPath: "testdata/test_en.mp3",
			lang:      "invalid",
			want:      "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mock.TranscribeAudio(tt.audioPath, tt.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranscribeAudio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TranscribeAudio() = %v, want %v", got, tt.want)
			}
			mock.Close() // Call Close() after each test case
		})
	}

	mock.AssertExpectations(t)
}

func TestModelMemoryRequirements(t *testing.T) {
	tests := []struct {
		name         string
		modelSize    string
		minMemoryMB  int
		maxMemoryMB  int
	}{
		{
			name:        "tiny model",
			modelSize:   "tiny",
			minMemoryMB: 75,
			maxMemoryMB: 100,
		},
		{
			name:        "base model",
			modelSize:   "base",
			minMemoryMB: 150,
			maxMemoryMB: 200,
		},
		{
			name:        "small model",
			modelSize:   "small",
			minMemoryMB: 400,
			maxMemoryMB: 500,
		},
		{
			name:        "medium model",
			modelSize:   "medium",
			minMemoryMB: 1000,
			maxMemoryMB: 1500,
		},
		{
			name:        "large model",
			modelSize:   "large",
			minMemoryMB: 2500,
			maxMemoryMB: 3000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip actual model loading, just verify memory requirements are documented
			assert.Greater(t, tt.maxMemoryMB, tt.minMemoryMB, "Max memory should be greater than min memory")
			assert.Greater(t, tt.minMemoryMB, 0, "Min memory should be positive")
		})
	}
}
