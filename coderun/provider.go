package coderun

import (
	"errors"
	"fmt"
	"path"

	"github.com/docker/docker/client"
)

type IProvider interface {
	Name() string
	Register(*RunEnvironment) bool
	ResourceRegister(*RunEnvironment)
	Resources() interface{}
	RegisteredResources() interface{}
	Setup(*RunEnvironment)
	Deploy(*RunEnvironment)
	Run(*RunEnvironment)
}

type ProviderHookFunc func(IProvider, *RunEnvironment)

type IResource interface {
	Name() string
	Register(*RunEnvironment) bool
}

type RunEnvironment struct {
	Name                string
	EntryPoint          string
	Providers           map[string]IProvider
	RegisteredProviders map[string]IProvider

	Cmd         []string
	CodeDir     string
	DependsDir  string
	IgnoreFiles []string
	Flags       map[string]*string
	CRDocker    ICRDocker
	Exec        func(...string) string
}

type IProviderEnv interface {
}

func CreateRunEnvironment() *RunEnvironment {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	cwd := Cwd()

	ignoreFiles := append(
		readIgnoreFile(path.Join(cwd, ".gitignore")),
		append(
			readIgnoreFile(path.Join(cwd, ".crignore")),
			".coderun",
		)...,
	)

	return &RunEnvironment{
		Providers: map[string]IProvider{
			"mount":  NewMountProvider(),
			"docker": NewDockerProvider(),
			//"lambda": NewAWSLambdaProvider(),
		},
		//Registered: map[string]map[*Provider]*Resource{},
		RegisteredProviders: map[string]IProvider{},
		CodeDir:             cwd,
		DependsDir:          "",
		IgnoreFiles:         ignoreFiles,
		Name:                path.Base(Cwd()),
		EntryPoint:          "lambda_handler",
		Flags:               make(map[string]*string),
		CRDocker:            &CRDocker{Client: cli},
		Exec:                Exec,
	}
}

func Setup(runEnv *RunEnvironment) (*RunEnvironment, error) {
	if p, _ := runEnv.Flags["provider"]; *p != "" {
		Logger.info.Printf("Registering provider %s", *p)
		runEnv.RegisteredProviders[*p] = runEnv.Providers[*p]
	} else {
		for n, p := range runEnv.Providers {
			//These probably should just be classes
			if p.Register(runEnv) {
				Logger.info.Printf("Registering provider %s", p.Name())
				runEnv.RegisteredProviders[n] = p
			}
		}
	}

	for _, p := range runEnv.RegisteredProviders {
		p.ResourceRegister(runEnv)
	}

	if len(runEnv.RegisteredProviders) <= 0 {
		return nil, errors.New(fmt.Sprintf("No providers found for this command"))
	}

	for _, p := range runEnv.RegisteredProviders {
		if p.Setup == nil {
			Logger.debug.Printf("No step `Setup` registered for provider `%s`", p.Name())
		} else {
			p.Setup(runEnv)
		}
	}
	for _, p := range runEnv.RegisteredProviders {
		if p.Deploy == nil {
			Logger.debug.Printf("No step `Deploy` registered for provider `%s`", p.Name())
		} else {
			p.Deploy(runEnv)
		}
	}
	for _, p := range runEnv.RegisteredProviders {
		if p.Run == nil {
			Logger.debug.Printf("No step `Run` registered for provider `%s`", p.Name())
		} else {
			p.Run(runEnv)
		}
	}

	return runEnv, nil
}
