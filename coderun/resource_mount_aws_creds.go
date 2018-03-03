package coderun

import (
	"io"
	"log"
	"regexp"
	"strings"
)

func NewAwsCredsMountResource() *AwsCredsMountResource {
	return &AwsCredsMountResource{}
}

type AwsCredsMountResource struct {
	fs *CoderunFs
}

func (cr *AwsCredsMountResource) Name() string { return "awsCreds" }

func (cr *AwsCredsMountResource) Register(e *RunEnvironment, p IProvider) bool {
	e.Registry.AddAt(SetupStep, &StepCallback{Step: "Setup", Provider: p, Resource: cr, Callback: cr.Setup})
	log.Printf("cr.fs: %v", cr.fs)
	return true
}

func (cr *AwsCredsMountResource) Path() string { return "/root/.aws" }

func (cr *AwsCredsMountResource) Fs() *CoderunFs { return cr.fs }

func (cr *AwsCredsMountResource) Setup(e *RunEnvironment, callback *StepCallback, currentStep *StepCallback) {
	Logger.debug.Printf("awsMountCreds setup")
	cr.fs = NewCoderunFs(cr.Path())
	cr.fs.Setup()
	cr.fs.AddFileResource(&credFile{})
	e.Registry.AddBefore( //Need to register this after the Fs is set up
		&StepSearch{Provider: regexp.MustCompile("docker"), Resource: regexp.MustCompile(".*"), Step: regexp.MustCompile("Run")},
		&StepCallback{Step: "ConnectDocker", Provider: callback.Provider, Callback: cr.fs.ConnectDocker})
}

type credFile struct {
	IFileResource
}

func (cf *credFile) Path() string { return "credentials" }

func (cf *credFile) Setup(e *RunEnvironment) { return }

func (cf *credFile) Open() io.Reader {
	return strings.NewReader("asdf")
}
