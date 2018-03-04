package coderun

type IDockerResource interface {
	IResource
	Setup(*StepCallback, *StepCallback)
	Run(*StepCallback, *StepCallback)
	RegisterMount(string, string)
}

type DockerResources map[string]IResource

func NewDockerProvider(r **RunEnvironment) IProvider {
	return &DockerProvider{
		resources: map[string]IDockerResource{
			"bash": NewBashResource(r),
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

func (p *DockerProvider) Register() bool {
	registered := false
	for name, r := range p.resources {
		if r.Register(p) {
			Logger.info.Printf("Registering resource %s", name)
			p.registeredResources[name] = r
			registered = true
		}
	}
	return registered
}
