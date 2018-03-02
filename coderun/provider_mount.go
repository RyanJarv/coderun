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
		resources: map[string]IMountResource{
			"awsCreds": NewAwsCredsMountResource(),
		},
		registeredResources: map[string]IMountResource{},
	}
}

type MountProvider struct {
	IProvider
	resources           map[string]IMountResource
	registeredResources map[string]IMountResource
	providerPreRunHook  map[string]ProviderHookFunc
}

func (p *MountProvider) Name() string {
	return "mount"
}

func (p *MountProvider) Register(e *RunEnvironment) bool {
	return true
}

func (p *MountProvider) Trigger(e *RunEnvironment) {
	Logger.info.Printf("Running provider %s", p.Name())
	p.Setup(e)
}

func (p *MountProvider) ResourceRegister(e *RunEnvironment) {
	for name, r := range p.resources {
		if r.Register(e) {
			Logger.info.Printf("Registering resource %s", name)
			p.registeredResources[name] = r
		}
	}

	if len(p.registeredResources) < 1 {
		log.Fatalf("Didn't find any registered lambda resources")
	}
}

func (p *MountProvider) Resources() interface{} {
	return p.resources
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
		p.connectDocker(r.Path(), r.Fs().LocalPath)
	}
}

func (p *MountProvider) connectDocker(runEnv *RunEnvironment, remotePath string, localPath string) {
	docker, ok := runEnv.RegisteredProviders["docker"]
	if ok != true {
		log.Info("Docker resource is not registered, will not set up shares")
	}

}
