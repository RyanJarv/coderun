package coderun

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Provider interface {
	Setup(*RunEnvironment)
	Run(*RunEnvironment)
}

type ProviderFactory func(conf *ProviderConfig) (Provider, error)

var providerFactories = make(map[string]ProviderFactory)

func Register(name string, factory ProviderFactory) {
	if factory == nil {
		log.Panicf("Provider %s does not exist.", name)
	}
	_, registered := providerFactories[name]
	if registered {
		log.Fatalf("Provider %s already registered. Ignoring.", name)
	}
	providerFactories[name] = factory
}

func init() {
	Register(".py", NewPythonProvider)
	Register(".rb", NewRubyProvider)
	Register(".go", NewGoProvider)
	Register(".js", NewJsProvider)
	Register(".sh", NewBashProvider)
}

func CreateProvider(c *ProviderConfig) (Provider, error) {
	factory, ok := providerFactories[c.Extension]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availableProviders := make([]string, len(providerFactories))
		for k, _ := range providerFactories {
			availableProviders = append(availableProviders, k)
		}
		return nil, errors.New(fmt.Sprintf("Invalid run provider name. Must be one of: %s", strings.Join(availableProviders, ", ")))
	}

	// Run the factory with the configuration.
	return factory(c)
}
