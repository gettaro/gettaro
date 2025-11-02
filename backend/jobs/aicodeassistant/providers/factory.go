package providers

import (
	"fmt"
)

// Factory creates AI code assistant providers
type Factory struct {
	providers map[string]AICodeAssistantProvider
}

// NewFactory creates a new provider factory
func NewFactory(providers []AICodeAssistantProvider) *Factory {
	factory := &Factory{
		providers: make(map[string]AICodeAssistantProvider),
	}

	for _, provider := range providers {
		factory.providers[provider.Name()] = provider
	}

	return factory
}

// GetProvider retrieves a provider by name
func (f *Factory) GetProvider(name string) (AICodeAssistantProvider, error) {
	provider, exists := f.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}

// NewFactoryFromProviders is a convenience function to create a factory
func NewFactoryFromProviders(providers ...AICodeAssistantProvider) *Factory {
	return NewFactory(providers)
}
