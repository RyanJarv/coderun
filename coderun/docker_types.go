package coderun

import (
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type dockerRunConfig struct {
	Image     string
	Port      int
	Cmd       []string
	SourceDir string
	DestDir   string
	Mounts    mount.Mount
}

type ResourceConfig struct {
	Extension     string
	Cmd           string
	Args          []string
	FullCmdString string
}

type ICRDocker interface {
	Pull(string)
	Run(dockerRunConfig)
	newImageName() string
	getOrBuildImage(string, ...[]string) string
	getImageName() string
}

type IRunEnvironment interface {
	Cmd() []string
}

type RunEnvironment struct {
	IRunEnvironment
	Cwd          string
	DockerClient *client.Client
	cmd          []string
	CRDocker     ICRDocker
}

func (d RunEnvironment) Cmd() []string {
	return d.cmd
}

//func (d RunEnvironment) CRDocker() *ICRDocker {
//	return &d.crdocker
//}

type RegisterOnCmdFunc func(cmd ...string) bool
type RunFunc func(IRunEnvironment)
type SetupFunc func(IRunEnvironment)

type IResource interface{}
type Resource struct {
	IResource
	RegisterOnCmd RegisterOnCmdFunc
	Setup         SetupFunc
	Run           RunFunc
}
