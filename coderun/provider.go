package coderun

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/docker/docker/client"
)

var pathResources = make(map[string]*Resource)

func Register(name string, provider *Resource) {
	if provider == nil {
		log.Panicf("Resource %s does not exist.", name)
	}
	_, registered := pathResources[name]
	if registered {
		log.Fatalf("Resource %s already registered. Ignoring.", name)
	}
	pathResources[name] = provider
}

func init() {
	Register("python", PythonResource())
	Register("ruby", RubyResource())
	Register("go", GoResource())
	Register("nodejs", JsResource())
	Register("bash", BashResource())
	Register("rails", RailsResource())
}

// func CreateRunEnvironment(jVj)(
// 	runEnv := &RunEnvironment{
// 		CRDocker:     &CRDocker{Client: cli},
// 		DockerClient: cli,
// 		Cmd:          flag.Args(),
// 	}
// )

func GetResource(c *ResourceConfig) (*Resource, error) {
	var provider *Resource
	for _, p := range pathResources {
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
		cmd:          flag.Args(),
	}

	provider.Setup(runEnv)
	provider.Run(runEnv)

	return provider, nil
}
