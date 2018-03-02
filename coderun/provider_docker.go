package coderun

import (
	"log"
)

type IDockerResource interface {
	IResource
	Setup(*RunEnvironment)
	Run(*RunEnvironment)
}

type DockerResources map[string]IResource

func NewDockerProvider() IProvider {
	return &DockerProvider{
		resources: map[string]IDockerResource{
			"bash": NewBashResource(),
		},
		registeredResources: map[string]IDockerResource{},
	}
}

type DockerProvider struct {
	IProvider
	resources           map[string]IDockerResource
	registeredResources map[string]IDockerResource
}

func (p *DockerProvider) Name() string {
	return "docker"
}

func (p *DockerProvider) Register(e *RunEnvironment) bool {
	return true
}

func (p *DockerProvider) Resources() interface{} {
	return p.resources
}

func (p *DockerProvider) RegisteredResources() interface{} {
	return p.registeredResources
}

func (p *DockerProvider) Trigger(e *RunEnvironment) {
	Logger.info.Printf("Running step `Run` for provider %s", p.Name())
	p.Run(e)
}

func (p *DockerProvider) ResourceRegister(e *RunEnvironment) {
	for name, r := range p.resources {
		if r.Register(e) {
			Logger.info.Printf("Registering resource %s", name)
			p.registeredResources[name] = r
		}
	}
	if len(p.registeredResources) < 1 {
		log.Fatalf("Didn't find any registered docker resources")
	}
}

func (p *DockerProvider) Setup(e *RunEnvironment) {
	for _, resource := range p.registeredResources {
		if resource.Setup == nil {
			Logger.info.Printf("No step Setup found for resource %s", resource.Name())
		} else {
			resource.Setup(e)
		}
	}
}

//func (p *DockerProvider) Deploy

func (p *DockerProvider) Run(e *RunEnvironment) {
	for _, resource := range p.registeredResources {
		if resource.Run == nil {
			Logger.info.Printf("No step Run found for resource %s", resource.Name())
		} else {
			resource.Run(e)
		}
	}
}
