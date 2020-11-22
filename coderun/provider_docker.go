package coderun

import "github.com/RyanJarv/coderun/coderun/lib"

type IDockerResource interface {
	IResource
}

type DockerResources map[string]IResource

func NewDockerProvider(r IRunEnvironment) IProvider {
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
	buildkit            *lib.CRDocker
}

func (p *DockerProvider) Name() string {
	return "docker"
}

func (p *DockerProvider) Register(e IRunEnvironment) bool {
	registered := false
	e.Registry().AddAt(TeardownStep+10, &StepCallback{Step: "Teardown", Provider: p, Callback: p.Teardown})
	for name, r := range p.resources {
		if r.Register(e, p) {
			Logger.Info.Printf("Registering resource %s", name)
			p.registeredResources[name] = r
			registered = true
		}
	}
	if registered == true {
		// This can be removed when buildkit get's merged into docker
		//(*p.env).Registry.AddAt(SetupStep-10, &StepCallback{Step: "Setup", Provider: p, Callback: p.Setup})
		//(*p.env).Registry.AddAt(TeardownStep, &StepCallback{Step: "Teardown", Provider: p, Callback: p.Teardown})
	}
	return registered
}

func (p *DockerProvider) Setup(callback *StepCallback, currentStep *StepCallback) {
	p.buildkit = lib.NewCRDocker()
	p.buildkit.Run(DockerRunConfig{
		Image:      "tonistiigi/buildkit",
		Attach:     false,
		Privileged: true,
		Port:       27467,
	})
}

func (p *DockerProvider) Teardown(callback *StepCallback, currentStep *StepCallback) {
	lib.NewCRDocker().DockerKillLabel("coderun")
}
