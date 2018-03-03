package coderun

import (
	"io"
	"log"
	"regexp"
)

type IMountResource interface {
	Name() string
	Register(*RunEnvironment, IProvider) bool
	Setup(*RunEnvironment, *StepCallback)
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
		registeredResources: []IMountResource{},
	}
}

type MountProvider struct {
	resources           []IMountResource
	registeredResources []IMountResource
}

func (p *MountProvider) Name() string {
	return "mount"
}

func (p *MountProvider) Register(e *RunEnvironment) bool {
	registered := false
	for _, r := range p.registeredResources {
		if r.Register(e, p) == true {
			registered = true
			e.Registry.AddBefore(
				&StepSearch{Provider: regexp.MustCompile("docker"), Step: regexp.MustCompile(".*"), Resource: regexp.MustCompile(".*")},
				&StepCallback{Step: "Setup", Provider: p, Callback: p.connectDocker})
		}
	}
	return registered
}

func (p *MountProvider) connectDocker(runEnv *RunEnvironment, c *StepCallback) {
	docker, ok := runEnv.RegisteredProviders["docker"]
	if ok != true {
		Logger.info.Printf("Docker resource is not registered, will not set up shares", docker)
	}

}
