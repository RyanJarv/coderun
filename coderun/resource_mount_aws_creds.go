package coderun

import (
	"io"
	"strings"
)

func NewAwsCredsMountResource() *AwsCredsMountResource {
	return &AwsCredsMountResource{}
}

type AwsCredsMountResource struct {
	IMountResource
	fs *CoderunFs
}

func (cr *AwsCredsMountResource) Name() string { return "awsCreds" }

func (cr *AwsCredsMountResource) Path() string { return "~/.aws" }

func (cr *AwsCredsMountResource) Register(r *RunEnvironment) bool { return true }

func (cr *AwsCredsMountResource) Fs() *CoderunFs { return cr.fs }

func (cr *AwsCredsMountResource) Setup(r *RunEnvironment) {
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
