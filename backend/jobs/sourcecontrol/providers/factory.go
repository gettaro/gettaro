package providers

import (
	"fmt"
)

// ProviderFactory creates a new provider instance based on the provider name
type ProviderFactory interface {
	GetProvider(providerName string) (SourceControlProvider, error)
}

type Factory struct {
	providers []SourceControlProvider
}

func NewFactory(providers []SourceControlProvider) ProviderFactory {
	return &Factory{
		providers: providers,
	}
}

func (f *Factory) GetProvider(providerName string) (SourceControlProvider, error) {
	for _, provider := range f.providers {
		if provider.Name() == providerName {
			return provider, nil
		}
	}
	return nil, fmt.Errorf("provider %s not found", providerName)
}
