package coderun

import (
	L "github.com/RyanJarv/coderun/coderun/logger"
)

func NewDockerProvider(r IRunEnvironment) IProvider {
	return &DockerProvider{
		resources: map[string]IProvider{
			"bash": NewBashResource(r),
		},
		registeredResources: map[string]IProvider{},
	}
}

type DockerProvider struct {
	IProvider
	resources           map[string]IProvider
	registeredResources map[string]IProvider
	buildkit            *CRDocker
}

func (d *DockerProvider) Name() string {
	return "docker"
}

func (d *DockerProvider) Register(e IRunEnvironment, p IProvider) bool {
	registered := false
	e.Registry().AddAt(TeardownStep+10, &StepCallback{Step: "Teardown", Provider: d, Callback: d.Teardown})
	for name, r := range d.resources {
		if r.Register(e, d) {
			L.Info.Printf("Registering resource %s", name)
			d.registeredResources[name] = r
			registered = true
		}
	}
	return registered
}

func (d *DockerProvider) Setup(callback *StepCallback, currentStep *StepCallback) {
	d.buildkit = NewCRDocker()
	d.buildkit.Run(DockerRunConfig{
		Image:      "tonistiigi/buildkit",
		Attach:     false,
		Privileged: true,
		Port:       27467,
	})
}

func (d *DockerProvider) Teardown(callback *StepCallback, currentStep *StepCallback) {
	NewCRDocker().DockerKillLabel("coderun")
}
