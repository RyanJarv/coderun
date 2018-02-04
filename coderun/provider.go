package coderun

import (
	"errors"
	"fmt"
	"log"
)

type ProviderDefault struct {
}

func (p *ProviderDefault) RegisterOnCmd(cmd string, args ...string) bool {
	return false
}

type Provider interface {
	RegisterOnCmd(cmd string, args ...string) bool
	Setup(*RunEnvironment)
	Run(*RunEnvironment)
}

var pathProviders = make(map[string]Provider)

func Register(name string, provider Provider) {
	if provider == nil {
		log.Panicf("Provider %s does not exist.", name)
	}
	_, registered := pathProviders[name]
	if registered {
		log.Fatalf("Provider %s already registered. Ignoring.", name)
	}
	pathProviders[name] = provider
}

func init() {
	Register("python", &PythonProvider{})
	Register("ruby", &RubyProvider{})
	Register("go", &GoProvider{})
	Register("nodejs", &JsProvider{})
	Register("bash", &BashProvider{})
	Register("rails", &RailsProvider{})
}

func CreateProvider(c *ProviderConfig) (Provider, error) {
	var provider Provider
	for _, p := range pathProviders {
		if p.RegisterOnCmd(c.Cmd, c.Args...) {
			provider = p
			break
		}
	}
	if provider == nil {
		return nil, errors.New(fmt.Sprintf("No providers found for this command"))
	}

	// Run the factory with the configuration.
	return provider, nil
}
