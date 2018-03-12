package coderun

import (
	"io"
	"net"
	"regexp"

	dclient "github.com/RyanJarv/dockersnitch/dockersnitch/client"
)

func NewSnitchDockerResource(r **RunEnvironment) *SnitchDockerResource {
	return &SnitchDockerResource{env: r, socket: "/var/run/dockersnitch.sock"}
}

type SnitchDockerResource struct {
	env     **RunEnvironment
	socket  string
	stream  net.Conn
	prompts chan string
	stdin   io.ReadWriteCloser
	stdout  io.ReadWriteCloser
}

func (sd *SnitchDockerResource) Name() string { return "snitchDocker" }

func (sd *SnitchDockerResource) Register(p IProvider) bool {
	(*sd.env).Registry.AddBefore(
		&StepSearch{Provider: regexp.MustCompile("docker"), Resource: regexp.MustCompile(".*"), Step: regexp.MustCompile("Setup")},
		&StepCallback{Step: "Setup", Provider: p, Resource: sd, Callback: sd.Setup},
	)
	return true
}

func (sd *SnitchDockerResource) Setup(callback *StepCallback, currentStep *StepCallback) {
	go (*sd.env).CRDocker.Run(dockerRunConfig{
		Image:      "dockersnitch",
		Privileged: true,
		Port:       33504,
		HostPort:   33505,
		PidMode:    "host",
	})

	dclient.Client("tcp", "localhost:33505")
}
