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
	Cwd          string
	CRDocker     *CRDocker
	DockerClient *client.Client
	Cmd          []string
}
