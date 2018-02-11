package coderun

import (
	"errors"
	"fmt"
	"path"
)

type Provider struct {
	Name                string
	Register            func(RunEnvironment) bool
	ResourceRegister    func(Provider, RunEnvironment)
	Setup               func(Provider, RunEnvironment) IProviderEnv
	Deploy              func(Provider, RunEnvironment, IProviderEnv)
	Run                 func(Provider, RunEnvironment, IProviderEnv)
	RegisteredResources map[string]*Resource
	Resources           map[string]*Resource
	ProviderEnv         IProviderEnv
}

type RunEnvironment struct {
	Name                string
	EntryPoint          string
	Providers           map[string]*Provider
	RegisteredProviders map[string]*Provider

	Cmd     []string
	CodeDir string
	Flags   map[string]*string
}

type IProviderEnv interface {
}

func CreateRunEnvironment() *RunEnvironment {
	return &RunEnvironment{
		Providers: map[string]*Provider{
			"docker": DockerProvider(),
			"lambda": AWSLambdaProvider(),
		},
		//Registered: map[string]map[*Provider]*Resource{},
		RegisteredProviders: map[string]*Provider{},
		CodeDir:             Cwd(),
		Name:                path.Base(Cwd()),
		EntryPoint:          "lambda_handler",
		Flags:               make(map[string]*string),
	}
}

func Setup(runEnv *RunEnvironment) (*RunEnvironment, error) {
	if p, _ := runEnv.Flags["provider"]; *p != "" {
		runEnv.RegisteredProviders[*p] = runEnv.Providers[*p]
	} else {
		for n, p := range runEnv.Providers {
			//These probably should just be classes
			if p.Register(*runEnv) {
				runEnv.RegisteredProviders[n] = p
			}
		}
	}

	for _, p := range runEnv.RegisteredProviders {
		p.ResourceRegister(*p, *runEnv)
	}

	if len(runEnv.RegisteredProviders) <= 0 {
		return nil, errors.New(fmt.Sprintf("No providers found for this command"))
	}

	for _, provider := range runEnv.RegisteredProviders {
		providerEnv := provider.Setup(*provider, *runEnv)

		runProviderStep("deploy", provider, *runEnv, providerEnv)
		runProviderStep("run", provider, *runEnv, providerEnv)
	}

	return runEnv, nil
}
