package transcribe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/giantswarm/mnote/internal/config"
)

// Transcriber interface defines the contract for audio transcription
type Transcriber interface {
	TranscribeAudio(audioPath, language string) (*TranscriptionResult, error)
}

// TranscriberImpl implements the Transcriber interface
type TranscriberImpl struct {
	config *config.Config
	client HTTPClient
}

// HTTPClient interface for mocking in tests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewTranscriber creates a new Transcriber instance
func NewTranscriber(cfg *config.Config) Transcriber {
	return &TranscriberImpl{
		config: cfg,
		client: &http.Client{},
	}
}

// TranscriptionResult represents the JSON response from the API
type TranscriptionResult struct {
	Text string `json:"text"`
}

// TranscribeAudio transcribes the audio file at the given path
func (t *TranscriberImpl) TranscribeAudio(audioPath, language string) (*TranscriptionResult, error) {
	// Open the audio file
	file, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()

	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add audio file
	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	// Add language parameter if not auto
	if language != "auto" {
		if err := writer.WriteField("language", language); err != nil {
			return nil, fmt.Errorf("failed to add language field: %w", err)
		}
	}

	// Add model parameter after language
	model := t.config.GetWhisperModel(language)
	fmt.Printf("Transcribing using model: %s (language: %s)\n", model, language)
	if err := writer.WriteField("model", model); err != nil {
		return nil, fmt.Errorf("failed to add model field: %w", err)
	}

	// Close multipart writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Get API URL from environment or config
	apiURL := os.Getenv("TRANSCRIPTION_API_URL")
	if apiURL == "" {
		apiURL = t.config.TranscriptionAPIURL
	}

	// Create request
	req, err := http.NewRequest("POST", apiURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result TranscriptionResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// MockHTTPClient implements HTTPClient for testing
type MockHTTPClient struct{}

func (m *MockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	// Mock response for testing
	jsonResponse := `{"text": "This is a mock transcription for testing purposes."}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(jsonResponse)),
	}, nil
}
