package coderun

import (
	"github.com/docker/docker/client"
)

type ProviderConfig struct {
	Extension     string
	Cmd           string
	Args          []string
	FullCmdString string
}

type RunEnvironment struct {
	Cwd           string
	DockerClient  *client.Client
	Cmd           []string
	ArgsString    string
	FullCmdString string
}
