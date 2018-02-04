package coderun

import (
	"fmt"
	"os"
)

type RubyProvider struct {
	ProviderDefault
	name string
}

func (p *RubyProvider) Setup(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "pull", "ruby:2.3")

	if _, err := os.Stat("./Gemfile"); err == nil {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "ruby:2.3", "bundler", "install", "--path", ".coderun/vendor/bundle")
	}
}

func (p *RubyProvider) Run(r *RunEnvironment) {
	name := newImageName()
	cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", name, "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "ruby:2.3", "ruby", r.ArgsString)
	cmd("/usr/local/bin/docker", "stop", name)
}
