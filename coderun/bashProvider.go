package coderun

import (
	"log"
	"regexp"
	"strings"
)

type BashProvider struct {
	ProviderDefault
	name string
}

func (p *BashProvider) RegisterOnCmd(cmd string, args ...string) bool {
	match, err := regexp.MatchString(`^(([^ ]+/)?bash .*|[\S]+\.sh)$`, strings.Join(append([]string{cmd}, args...), " "))
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func (p *BashProvider) Setup(r *RunEnvironment) {
	dockerPull(r.DockerClient, "bash")
}

func (p *BashProvider) Run(r *RunEnvironment) {
	dockerRun(dockerRunConfig{Client: r.DockerClient, Image: "ubuntu", DestDir: "/usr/src/myapp", SourceDir: r.Cwd, Cmd: append([]string{"bash"}, r.Cmd...)})
}
