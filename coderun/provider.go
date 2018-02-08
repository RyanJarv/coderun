package coderun

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/docker/docker/client"
)

var pathProviders = make(map[string]*Provider)

type RegisterOnCmdFunc func(cmd ...string) bool
type SetupFunc func(*RunEnvironment)
type RunFunc func(*RunEnvironment)

type Provider struct {
	RegisterOnCmd RegisterOnCmdFunc
	Setup         SetupFunc
	Run           RunFunc
}

func Register(name string, provider *Provider) {
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
	Register("python", PythonProvider())
	Register("ruby", RubyProvider())
	Register("go", GoProvider())
	Register("nodejs", JsProvider())
	Register("bash", BashProvider())
	Register("rails", RailsProvider())
}

func GetProvider(c *ProviderConfig) (*Provider, error) {
	var provider *Provider
	for _, p := range pathProviders {
		if p.RegisterOnCmd(append([]string{c.Cmd}, c.Args...)...) {
			provider = p
			break
		}
	}
	if provider == nil {
		return nil, errors.New(fmt.Sprintf("No providers found for this command"))
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	runEnv := &RunEnvironment{
		CRDocker:     &CRDocker{Client: cli},
		DockerClient: cli,
		Cmd:          flag.Args(),
	}

	provider.Setup(runEnv)
	provider.Run(runEnv)

	return provider, nil
}
