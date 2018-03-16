package coderun

import (
	"path"
)

type IProvider interface {
	Name() string
	Register(IRunEnvironment) bool
}

type ProviderHookFunc func(IProvider, IRunEnvironment)

type IResource interface {
	Name() string
	Register(IRunEnvironment, IProvider) bool
}

type IRunEnvironment interface {
	Providers() map[string]IProvider
	Registry() *Registry
	Cmd() []string
}

type RunEnvironment struct {
	Name                string
	EntryPoint          string
	providers           map[string]IProvider
	registeredProviders map[string]IProvider

	cmd         []string
	CodeDir     string
	DependsDir  string
	IgnoreFiles []string
	Flags       map[string]*string
	CRDocker    ICRDocker
	Exec        func(...string) string
	registry    *Registry
}

func (e *RunEnvironment) Providers() map[string]IProvider { return e.providers }
func (e *RunEnvironment) Cmd() []string                   { return e.cmd }
func (e *RunEnvironment) Registry() *Registry             { return e.registry }

type IProviderEnv interface {
}

func CreateRunEnvironment() *RunEnvironment {
	cwd := Cwd()

	ignoreFiles := append(
		readIgnoreFile(path.Join(cwd, ".gitignore")),
		append(
			readIgnoreFile(path.Join(cwd, ".crignore")),
			".coderun",
		)...,
	)

	var runEnv *RunEnvironment
	runEnv = &RunEnvironment{
		providers: map[string]IProvider{
			"mount":  NewMountProvider(runEnv),
			"docker": NewDockerProvider(runEnv),
			"shell":  NewShellProvider(runEnv),
			"snitch": NewSnitchProvider(runEnv),
			//"lambda": NewAWSLambdaProvider(),
		},
		//Registered: map[string]map[*Provider]*Resource{},
		registeredProviders: map[string]IProvider{},
		CodeDir:             cwd,
		DependsDir:          "",
		IgnoreFiles:         ignoreFiles,
		Name:                path.Base(Cwd()),
		EntryPoint:          "lambda_handler",
		Flags:               make(map[string]*string),
		Exec:                Exec,
		registry:            NewRegistry(),
	}
	return runEnv
}

func Setup(e *RunEnvironment, cmd []string) (IRunEnvironment, error) {
	e.cmd = cmd
	for _, p := range (*e).providers {
		p.Register(e)
	}

	e.registry.Run()
	Logger.info.Printf("Done running steps")

	return e, nil
}
