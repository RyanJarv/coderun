package coderun

import (
	"log"

	"github.com/docker/docker/client"
)

type dockerProviderEnv struct {
	IProviderEnv
	CRDocker ICRDocker
	cli      *client.Client
	Exec     func(...string) string
}

func DockerProvider() *Provider {
	return &Provider{
		Name:             "docker",
		Register:         dockerRegister,
		ResourceRegister: dockerResourceRegister,
		Setup:            dockerSetup,
		Run:              dockerRun,
		Resources: map[string]*Resource{
			"python": PythonResource(),
			"ruby":   RubyResource(),
			"go":     GoResource(),
			"nodejs": JsResource(),
			"bash":   BashResource(),
			"rails":  RailsResource(),
		},
		RegisteredResources: map[string]*Resource{},
		ProviderEnv:         nil,
	}
}

func dockerRegister(r RunEnvironment) bool {
	return true
}

func dockerResourceRegister(p Provider, runEnv RunEnvironment) {
	for n, r := range p.Resources {
		if r.Register(runEnv, p) {
			p.RegisteredResources[n] = r
		}
	}
	if len(p.RegisteredResources) != 1 {
		log.Fatalf("Found %d registered docker resources, should be one.\n", len(p.RegisteredResources))
	}
}

func dockerSetup(provider Provider, r RunEnvironment) IProviderEnv {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	providerEnv := dockerProviderEnv{
		CRDocker: &CRDocker{Client: cli},
		cli:      cli,
		Exec:     Exec,
	}
	for _, resource := range provider.RegisteredResources {
		resource.Setup(r, providerEnv)
	}
	return providerEnv
}

func dockerRun(provider Provider, r RunEnvironment, p IProviderEnv) {
	providerEnv := p.(dockerProviderEnv)
	for _, resource := range provider.RegisteredResources {
		resource.Run(r, providerEnv)
	}
}
