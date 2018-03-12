package coderun

import (
	"io"
)

type IMountResource interface {
	Name() string
	Register(IProvider) bool
	Setup(*StepCallback, *StepCallback)
	Path() string
	Fs() CoderunFs
}

type IFileResource interface {
	Path() string
	Setup()
	Open() io.Reader
}

func NewMountProvider(r **RunEnvironment) IProvider {
	return &MountProvider{
		resources: []IMountResource{
			NewAwsCredsMountResource(r),
		},
	}
}

type MountProvider struct {
	resources []IMountResource
}

func (p *MountProvider) Name() string {
	return "mount"
}

func (p *MountProvider) Register() bool {
	registered := false
	for _, r := range p.resources {
		if r.Register(p) == true {
			registered = true
		}
	}
	return registered
}
