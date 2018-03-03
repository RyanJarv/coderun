package coderun

import (
	"io"
	"log"
)

type IMountResource interface {
	IResource
	Setup(*RunEnvironment)
	Path() string
	Fs() *CoderunFs
}

type IFileResource interface {
	IResource
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
	IProvider
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
		}
	}
	return registered
}

func (p *MountProvider) Trigger(e *RunEnvironment) {
	Logger.info.Printf("Running provider %s", p.Name())
	p.Setup(e)
}

func (p *MountProvider) ResourceRegister(e *RunEnvironment) {
	for name, r := range p.resources {
		if r.Register(e, p) == true {
			Logger.info.Printf("Registering resource %s", name)
			p.registeredResources[name] = r
		}
	}

	if len(p.registeredResources) < 1 {
		log.Fatalf("Didn't find any registered lambda resources")
	}
}

func (p *MountProvider) Setup(e *RunEnvironment) {
	for _, r := range p.registeredResources {
		if r.Setup == nil {
			Logger.info.Printf("No step Setup found for resource %s", r.Name)
		} else {
			r.Setup(e)
		}
	}
	for _, r := range p.registeredResources {
		p.connectDocker(e, r.Path(), r.Fs().LocalPath)
	}
}

func (p *MountProvider) connectDocker(runEnv *RunEnvironment, remotePath string, localPath string) {
	docker, ok := runEnv.RegisteredProviders["docker"]
	if ok != true {
		Logger.info.Printf("Docker resource is not registered, will not set up shares", docker)
	}

}
