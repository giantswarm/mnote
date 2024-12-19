package interfaces

// Transcriber defines the interface for audio transcription services
type Transcriber interface {
	// TranscribeAudio transcribes an audio file to text
	// lang can be a specific language code or "auto" for auto-detection
	TranscribeAudio(audioPath string, lang string) (string, error)
}
