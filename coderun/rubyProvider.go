package coderun

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type RubyProvider struct {
	ProviderDefault
	name string
}

func (p *RubyProvider) RegisterOnCmd(cmd string, args ...string) bool {
	match, err := regexp.MatchString(`^(([^ ]+/)?ruby[0-9.]* .*|[\S]+\.rb)$`, strings.Join(append([]string{cmd}, args...), " "))
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func (p *RubyProvider) Setup(r *RunEnvironment) {
	dockerPull(r.DockerClient, "ruby:2.3")

	if _, err := os.Stat("./Gemfile"); err == nil {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "ruby:2.3", "bundler", "install", "--path", ".coderun/vendor/bundle")
	}
}

func (p *RubyProvider) Run(r *RunEnvironment) {
	dockerRun(dockerRunConfig{Client: r.DockerClient, Image: "ruby:2.3", DestDir: "/usr/src/myapp", SourceDir: r.Cwd, Cmd: append([]string{"ruby"}, r.Cmd...)})
}
