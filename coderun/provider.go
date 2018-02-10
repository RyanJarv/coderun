package coderun

import (
	"errors"
	"flag"
	"fmt"
)

type Provider struct {
	Name             string
	Register         func(RunEnvironment) bool
	ResourceRegister func(*Resource, RunEnvironment, IProviderEnv) bool
	Setup            func(RunEnvironment) IProviderEnv
	Run              func(RunEnvironment, IProviderEnv)
	Resources        map[string]*Resource
	ProviderEnv      IProviderEnv
}

type RunEnvironment struct {
	Providers  map[string]*Provider
	Registered map[string]*Resource

	Cmd []string
	Cwd string
}

type IProviderEnv interface {
}

func GetRunEnvironment() *RunEnvironment {
	return &RunEnvironment{
		Providers: map[string]*Provider{
			"docker": DockerProvider(),
			//"s3":   S3Resource(),
		},
		Registered: map[string]*Resource{},
		Cmd:        flag.Args(),
	}
}

type ResourceConfig struct {
	Cmd string
}

func Setup(c *ResourceConfig) (*RunEnvironment, error) {
	runEnv := GetRunEnvironment()

	registered := map[*Provider]*Resource{}
	for _, provider := range runEnv.Providers {
		if provider.Register(*runEnv) {
			for _, resource := range provider.Resources {
				if provider.ResourceRegister(resource, *runEnv, provider.ProviderEnv) {
					registered[provider] = resource
				}
			}
		}
	}
	if len(registered) <= 0 {
		return nil, errors.New(fmt.Sprintf("No providers found for this command"))
	}

	for provider, resource := range registered {
		provider.ProviderEnv = provider.Setup(*runEnv)
		resource.Setup(*runEnv, provider.ProviderEnv)
		provider.Run(*runEnv, provider.ProviderEnv)
		resource.Run(*runEnv, provider.ProviderEnv)
	}

	return runEnv, nil
}
