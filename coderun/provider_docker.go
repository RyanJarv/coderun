package coderun

import (
	"log"
)

type dockerProviderEnv struct {
	IProviderEnv
}

func DockerProvider() *Provider {
	return &Provider{
		Name:             "docker",
		Register:         dockerRegister,
		ResourceRegister: dockerResourceRegister,
		Setup:            dockerSetup,
		Run:              dockerRun,
		Resources: map[string]*Resource{
			"bundler": BundlerResource(),
			"ruby":    RubyResource(),
			"bash":    BashResource(),
		},
		RegisteredResources: map[string]*Resource{},
		ProviderEnv:         nil,
	}
}

func dockerRegister(r *RunEnvironment) bool {
	return true
}

func dockerResourceRegister(p Provider, runEnv *RunEnvironment) {
	for name, r := range p.Resources {
		if r.Register(runEnv, p) {
			Logger.info.Printf("Registering resource %s", name)
			p.RegisteredResources[name] = r
		}
	}
	if len(p.RegisteredResources) < 1 {
		log.Fatalf("Didn't find any registered docker resources")
	}
}

func dockerSetup(provider Provider, r *RunEnvironment) IProviderEnv {
	providerEnv := dockerProviderEnv{}

	for _, resource := range provider.RegisteredResources {
		if resource.Setup == nil {
			Logger.info.Printf("No step Setup found for resource %s", resource.Name)
		} else {
			resource.Setup(r, providerEnv)
		}
	}
	return providerEnv
}

func dockerRun(provider Provider, r *RunEnvironment, p IProviderEnv) {
	providerEnv := r
	for _, resource := range provider.RegisteredResources {
		if resource.Run == nil {
			Logger.info.Printf("No step Run found for resource %s", resource.Name)
		} else {
			resource.Run(r, providerEnv)
		}
	}
}
