package coderun

import (
	"io"
	"strings"
)

func NewAwsCredsMountResource() *AwsCredsMountResource {
	return &AwsCredsMountResource{}
}

type AwsCredsMountResource struct {
	fs *CoderunFs
}

func (cr *AwsCredsMountResource) Name() string { return "awsCreds" }

func (cr *AwsCredsMountResource) Register(r *RunEnvironment, p IProvider) bool {
	r.Registry.AddAt(SetupStep, &StepCallback{Step: "Setup", Provider: p, Resource: cr, Callback: cr.Setup})
	return true
}

func (cr *AwsCredsMountResource) Path() string { return "~/.aws" }

func (cr *AwsCredsMountResource) Fs() *CoderunFs { return cr.fs }

func (cr *AwsCredsMountResource) Setup(r *RunEnvironment, c *StepCallback) {
	cr.fs = NewCoderunFs(cr.Path())
	cr.fs.AddFileResource(&credFile{})
}

type credFile struct {
	IFileResource
}

func (cf *credFile) Name() string { return "awsCreds" }
func (cf *credFile) Path() string { return "~/.aws/credentials" }

func (cf *credFile) Setup(e *RunEnvironment) { return }

func (cf *credFile) Open() io.Reader {
	return strings.NewReader("asdf")
}
