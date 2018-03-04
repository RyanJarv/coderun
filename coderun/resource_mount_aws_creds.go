package coderun

import (
	"io"
	"regexp"
	"strings"
)

func NewAwsCredsMountResource(r **RunEnvironment) *AwsCredsMountResource {
	return &AwsCredsMountResource{env: r}
}

type AwsCredsMountResource struct {
	env **RunEnvironment
	fs  CoderunFs
}

func (cr *AwsCredsMountResource) Name() string { return "awsCreds" }

func (cr *AwsCredsMountResource) Register(p IProvider) bool {
	(*cr.env).Registry.AddAt(SetupStep, &StepCallback{Step: "Setup", Provider: p, Resource: cr, Callback: cr.Setup})
	(*cr.env).Registry.AddAt(TeardownStep, &StepCallback{
		Step:     "Unmount",
		Provider: p,
		Callback: func(*StepCallback, *StepCallback) { cr.fs.server.Unmount() }})
	(*cr.env).Registry.AddBefore( //Need to register this after the Fs is set up
		&StepSearch{Provider: regexp.MustCompile("docker"), Resource: regexp.MustCompile(".*"), Step: regexp.MustCompile("Run")},
		&StepCallback{Step: "ConnectDocker", Provider: p, Callback: cr.fs.ConnectDocker})
	return true
}

func (cr *AwsCredsMountResource) Path() string { return "/root/.aws" }

func (cr *AwsCredsMountResource) Fs() CoderunFs { return cr.fs }

func (cr *AwsCredsMountResource) Setup(callback *StepCallback, currentStep *StepCallback) {
	Logger.debug.Printf("awsMountCreds setup")
	cr.fs = NewCoderunFs(cr.Path())
	cr.fs.AddFileResource(&credFile{})
	cr.fs.Setup()
	go cr.fs.Serve()
}

type credFile struct {
	IFileResource
}

func (cf *credFile) Path() string { return "credentials" }

func (cf *credFile) Setup() { return }

func (cf *credFile) Open() io.Reader {
	return strings.NewReader("asdf")
}
