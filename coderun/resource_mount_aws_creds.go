package coderun

import (
	"io"
	"regexp"
	"strings"
)

func NewAwsCredsMountResource(r IRunEnvironment) *AwsCredsMountResource {
	return &AwsCredsMountResource{}
}

type AwsCredsMountResource struct {
	fs  CoderunFs
	env IRunEnvironment
}

func (cr *AwsCredsMountResource) Name() string { return "awsCreds" }

func (cr *AwsCredsMountResource) Register(e IRunEnvironment, p IProvider) bool {
	cr.env = e
	if len(e.Cmd()) == 0 {
		return false
	}
	e.Registry().AddAt(SetupStep, &StepCallback{Step: "Setup", Provider: p, Resource: cr, Callback: cr.Setup})
	e.Registry().AddAt(TeardownStep, &StepCallback{
		Step:     "Unmount",
		Provider: p,
		Resource: cr,
		Callback: func(*StepCallback, *StepCallback) { cr.fs.server.Unmount() }})
	e.Registry().AddBefore( //Need to register this after the Fs is set up
		&StepSearch{Provider: regexp.MustCompile("docker"), Resource: regexp.MustCompile(".*"), Step: regexp.MustCompile("Run")},
		&StepCallback{Step: "ConnectDocker", Provider: p, Resource: cr, Callback: cr.fs.ConnectDocker})
	return true
}

func (cr *AwsCredsMountResource) Path() string { return "/root/.aws" }

func (cr *AwsCredsMountResource) Fs() CoderunFs { return cr.fs }

func (cr *AwsCredsMountResource) Setup(callback *StepCallback, currentStep *StepCallback) {
	Logger.debug.Printf("awsMountCreds setup")
	cr.fs = NewCoderunFs(cr.Path())
	cr.fs.AddFileResource(&credFile{env: cr.env})
	cr.fs.Setup()
	go cr.fs.Serve()
}

type credFile struct {
	IFileResource
	env IRunEnvironment
}

func (cf *credFile) Path() string { return "credentials" }

func (cf *credFile) Setup() { return }

func (cf *credFile) Open() io.Reader {
	resp := cf.env.Stdin().Prompt("!!! Script is attempting to read ~/.aws/credentials, is this expected? [yes/no] ")

	var out io.Reader
	if resp == "yes" {
		out = strings.NewReader("***super secret keys stored on the host machine***\n")
	} else {
		Logger.warn.Printf("Restricting access to ~/.aws/credentials")
		out = strings.NewReader("")
	}
	return out
}
