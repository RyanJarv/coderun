package coderun

import (
	"path"

	"github.com/docker/docker/client"
)

type IProvider interface {
	Name() string
	Register(*RunEnvironment) bool
}

type ProviderHookFunc func(IProvider, *RunEnvironment)

type IResource interface {
	Name() string
	Register(*RunEnvironment, IProvider) bool
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
		Registry:            NewRegistry(),
	}
}

func Setup(runEnv *RunEnvironment) (*RunEnvironment, error) {
	for _, p := range runEnv.Providers {
		p.Register(runEnv)
	}

	runEnv.Registry.Run(runEnv)
	Logger.info.Printf("Done running steps")

	return runEnv, nil
}
