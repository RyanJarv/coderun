package coderun

type IDockerResource interface {
	IResource
	Setup(*RunEnvironment, *StepCallback)
	Run(*RunEnvironment, *StepCallback)
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
	registered := false
	for name, r := range p.resources {
		if r.Register(e, p) {
			Logger.info.Printf("Registering resource %s", name)
			p.registeredResources[name] = r
			registered = true
		}
	}
	return registered
}
