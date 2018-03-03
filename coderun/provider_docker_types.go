package coderun

import (
	"github.com/docker/docker/api/types/mount"
)

type dockerRunConfig struct {
	Image     string
	Port      int
	Cmd       []string
	SourceDir string
	DestDir   string
	Mounts    mount.Mount
}

type ICRDocker interface {
	Pull(string)
	Run(dockerRunConfig)
	RegisterMount(string, string)
	newImageName() string
	getOrBuildImage(string, ...[]string) string
	getImageName() string
}
