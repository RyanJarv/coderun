package coderun

import (
	"io"
	"log"
)

type IMountResource interface {
	Name() string
	Register(*RunEnvironment, IProvider) bool
	Setup(*RunEnvironment, *StepCallback, *StepCallback)
	Path() string
	Fs() *CoderunFs
}

type IFileResource interface {
	Name() string
	Path() string
	Setup(*RunEnvironment)
	Open() io.Reader
}

func NewMountProvider() IProvider {
	log.Printf("Settin up MountProvider")
	return &MountProvider{
		resources: []IMountResource{
			NewAwsCredsMountResource(),
		},
	}
}

type MountProvider struct {
	resources []IMountResource
}

func (p *MountProvider) Name() string {
	return "mount"
}

func (p *MountProvider) Register(e *RunEnvironment) bool {
	registered := false
	for _, r := range p.resources {
		if r.Register(e, p) == true {
			registered = true
		}
	}
	return registered
}
