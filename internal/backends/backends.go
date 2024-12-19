package backends

import (
	"net/http"

	"github.com/giantswarm/mnote/internal/config"
	"github.com/giantswarm/mnote/internal/interfaces"
	"github.com/giantswarm/mnote/internal/registry"
	"github.com/giantswarm/mnote/internal/transcribe"
	"github.com/giantswarm/mnote/internal/whisper"
)

// RegisterBackends registers all available transcription backends
func RegisterBackends() {
	// Register KubeAI backend
	registry.RegisterBackend("kubeai", func(cfg *config.Config) (interfaces.Transcriber, error) {
		return &transcribe.KubeAITranscriber{
			Config: cfg,
			Client: &http.Client{},
		}, nil
	})

	// Register local whisper backend
	registry.RegisterBackend("local", func(cfg *config.Config) (interfaces.Transcriber, error) {
		return whisper.NewLocalWhisper(cfg)
	})
}
