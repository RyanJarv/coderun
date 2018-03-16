package coderun

import (
	"regexp"
	"time"

	dclient "github.com/RyanJarv/dockersnitch/dockersnitch/client"
)

func NewSnitchDockerResource(r IRunEnvironment) *SnitchDockerResource {
	return &SnitchDockerResource{}
}

type SnitchDockerResource struct {
	dockersnitch **CRDocker
}

func (sd *SnitchDockerResource) Name() string { return "snitchDocker" }

func (sd *SnitchDockerResource) Register(e IRunEnvironment, p IProvider) bool {
	e.Registry().AddBefore(
		&StepSearch{Provider: regexp.MustCompile("docker"), Resource: regexp.MustCompile(".*"), Step: regexp.MustCompile("Setup")},
		&StepCallback{Step: "Setup", Provider: p, Resource: sd, Callback: sd.Setup},
	)
	e.Registry().AddAt(TeardownStep, &StepCallback{Step: "Teardown", Provider: p, Resource: sd, Callback: sd.Teardown})
	return true
}

func (sd *SnitchDockerResource) Setup(callback *StepCallback, currentStep *StepCallback) {
	s := NewCRDocker()
	sd.dockersnitch = &s
	(*sd.dockersnitch).Run(dockerRunConfig{
		Image:      "dockersnitch",
		Attach:     false,
		Privileged: true,
		Port:       33504,
		HostPort:   33505,
		PidMode:    "host",
	})

	dclient.Client("tcp", "localhost:33505")
}

func (sd SnitchDockerResource) Teardown(callback *StepCallback, currentStep *StepCallback) {
	(*sd.dockersnitch).Teardown(4 * time.Second)
}
