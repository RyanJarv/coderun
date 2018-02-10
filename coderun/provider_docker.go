package coderun

import (
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
		ProviderEnv: nil,
	}
}

func dockerRegister(r RunEnvironment) bool {
	return MatchCommandOrExt(r.Cmd, "bash", ".sh")
}

func dockerResourceRegister(resource *Resource, runEnv RunEnvironment, p IProviderEnv) bool {
	return resource.Register(runEnv, p.(dockerProviderEnv))
}

func dockerSetup(r RunEnvironment) IProviderEnv {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	providerEnv := &dockerProviderEnv{
		CRDocker: &CRDocker{Client: cli},
		cli:      cli,
		Exec:     Exec,
	}
	return providerEnv
}

func dockerRun(r RunEnvironment, p IProviderEnv) {
}
