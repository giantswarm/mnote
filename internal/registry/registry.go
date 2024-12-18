package registry

import (
	"fmt"
	"sync"

	"github.com/giantswarm/mnote/internal/config"
	"github.com/giantswarm/mnote/internal/interfaces"
)

// BackendFactory is a function type that creates a new Transcriber instance
type BackendFactory func(cfg *config.Config) (interfaces.Transcriber, error)

var (
	// registry holds all registered transcription backends
	registry = make(map[string]BackendFactory)
	mu       sync.RWMutex
)

// RegisterBackend registers a new transcription backend with the given name
func RegisterBackend(name string, factory BackendFactory) {
	mu.Lock()
	defer mu.Unlock()
	registry[name] = factory
}

// GetBackend returns a new instance of the requested transcription backend
func GetBackend(name string, cfg *config.Config) (interfaces.Transcriber, error) {
	mu.RLock()
	factory, ok := registry[name]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown transcription backend: %s", name)
	}
	return factory(cfg)
}
