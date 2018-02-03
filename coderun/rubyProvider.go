package coderun

import (
	"fmt"
	"os"
)

func NewRubyProvider(conf *ProviderConfig) (Provider, error) {
	return &RubyProvider{
		name: "ruby",
	}, nil
}

type RubyProvider struct {
	name string
}

func (p *RubyProvider) Setup(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "pull", "ruby:2.1")

	if _, err := os.Stat("./Gemfile"); err == nil {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", "my-running-script1", "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "ruby:2.1", "bundler", "install", "--path", ".coderun/vendor/bundle")
	}
}

func (p *RubyProvider) Run(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "ruby:2.1", "ruby", r.FilePath)
}
