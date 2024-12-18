package backends

import (
	"fmt"
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
		modelName := whisper.GetDefaultModel(cfg.DefaultLanguage).Name
		model, ok := cfg.Catalog[modelName]
		if !ok {
			return nil, fmt.Errorf("model %s not found in catalog", modelName)
		}
		return whisper.NewLocalWhisper(model)
	})
}
