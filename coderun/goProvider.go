package coderun

import (
	"log"
	"os"
	"regexp"
	"strings"
)

type GoProvider struct {
	ProviderDefault
	name string
}

func (p *GoProvider) RegisterOnCmd(cmd string, args ...string) bool {
	match, err := regexp.MatchString(`^(([^ ]+/)?go .*|[\S]+\.go)$`, strings.Join(append([]string{cmd}, args...), " "))
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func (p *GoProvider) Setup(r *RunEnvironment) {
	dockerPull(r.DockerClient, "golang")

	if _, err := os.Stat("./glide.lock"); os.IsNotExist(err) {
		dockerRun(dockerRunConfig{Client: r.DockerClient, Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: r.Cwd, Cmd: []string{"sh", "-c", "curl https://glide.sh/get | sh && glide init --non-interactive && glide install"}})
	}
}

func (p *GoProvider) Run(r *RunEnvironment) {
	dockerRun(dockerRunConfig{Client: r.DockerClient, Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: r.Cwd, Cmd: []string{"go", "run"}})
}
