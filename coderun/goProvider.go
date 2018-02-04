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
	dockerPull("golang")

	if _, err := os.Stat("./glide.lock"); os.IsNotExist(err) {
		dockerRun("golang", 1234, "/go/src/localhost/myapp", "sh", "-c", "curl https://glide.sh/get | sh && glide init --non-interactive && glide install")
	}
}

func (p *GoProvider) Run(r *RunEnvironment) {
	dockerRun("golang", 1234, "/go/src/localhost/myapp", append([]string{"go", "run", r.Cmd}, r.Args...)...)
}
