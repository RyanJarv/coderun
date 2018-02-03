package coderun

import (
	"fmt"
)

func NewBashProvider(conf *ProviderConfig) (Provider, error) {
	return &BashProvider{
		name: "ruby",
	}, nil
}

type BashProvider struct {
	name string
}

func (p *BashProvider) Setup(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "pull", "ubuntu")
}

func (p *BashProvider) Run(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "ubuntu", "bash", r.FilePath)
}
