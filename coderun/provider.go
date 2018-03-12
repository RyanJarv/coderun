package coderun

import (
	"path"
)

type IProvider interface {
	Name() string
	Register() bool
}

type ProviderHookFunc func(IProvider, *RunEnvironment)

type IResource interface {
	Name() string
	Register(IProvider) bool
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
	Registry    *Registry
}

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
		Providers: map[string]IProvider{
			"mount":  NewMountProvider(&runEnv),
			"docker": NewDockerProvider(&runEnv),
			"snitch": NewSnitchProvider(&runEnv),
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
		Exec:                Exec,
		Registry:            NewRegistry(),
	}
	return runEnv
}

func Setup(runEnv *RunEnvironment) (*RunEnvironment, error) {
	for _, p := range runEnv.Providers {
		p.Register()
	}

	runEnv.Registry.Run()
	Logger.info.Printf("Done running steps")

	return runEnv, nil
}
