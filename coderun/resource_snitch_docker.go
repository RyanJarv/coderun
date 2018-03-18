package coderun

import (
	"regexp"
	"time"

	dsclient "github.com/RyanJarv/dockersnitch/dockersnitch/client"
)

func NewSnitchDockerResource(r IRunEnvironment) *SnitchDockerResource {
	return &SnitchDockerResource{}
}

type SnitchDockerResource struct {
	dockersnitch *CRDocker
	dsclient     *dsclient.Client
	env          IRunEnvironment
}

func (sd *SnitchDockerResource) Name() string { return "snitchDocker" }

func (sd *SnitchDockerResource) Register(e IRunEnvironment, p IProvider) bool {
	sd.env = e
	e.Registry().AddBefore(
		&StepSearch{Provider: regexp.MustCompile("docker"), Resource: regexp.MustCompile(".*"), Step: regexp.MustCompile("Setup")},
		&StepCallback{Step: "Setup", Provider: p, Resource: sd, Callback: sd.Setup},
	)
	e.Registry().AddAfter(
		&StepSearch{Provider: regexp.MustCompile("docker"), Resource: regexp.MustCompile(".*"), Step: regexp.MustCompile("Teardown")},
		&StepCallback{Step: "Teardown", Provider: p, Resource: sd, Callback: sd.Teardown},
	)
	return true
}

func (sd *SnitchDockerResource) Setup(callback *StepCallback, currentStep *StepCallback) {
	sd.dockersnitch = NewCRDocker()
	sd.dockersnitch.Run(dockerRunConfig{
		Image:      "dockersnitch",
		Attach:     false,
		Privileged: true,
		Port:       33504,
		HostPort:   33505,
		PidMode:    "host",
	})

	sd.dsclient = &dsclient.Client{Ask: sd.env.Ask}
	sd.dsclient.Start("tcp", "localhost:33505")
}

func (sd *SnitchDockerResource) Teardown(callback *StepCallback, currentStep *StepCallback) {
	sd.dockersnitch.Teardown(4 * time.Second)
	sd.dsclient.Teardown()
}
