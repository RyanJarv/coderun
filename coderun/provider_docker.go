package coderun

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

func (p *DockerProvider) Trigger(e *RunEnvironment) {
	Logger.info.Printf("Running step `Run` for provider %s", p.Name())
	p.Run(e)
}
