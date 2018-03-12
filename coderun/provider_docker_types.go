package coderun

import (
	"io"

	"github.com/docker/docker/api/types/mount"
)

type dockerRunConfig struct {
	Image       string
	Port        int
	HostPort    int
	Cmd         []string
	SourceDir   string
	DestDir     string
	Mounts      mount.Mount
	Privileged  bool
	NetworkMode string
	PidMode     string
	Stdin       io.Reader
	Stdout      io.Writer
}

type ICRDocker interface {
	Pull(string)
	Run(dockerRunConfig)
	RegisterMount(string, string)
	newImageName() string
	getOrBuildImage(string, ...[]string) string
	getImageName() string
}
