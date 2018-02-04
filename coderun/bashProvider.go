package coderun

import (
	"fmt"
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
	cmd("/usr/local/bin/docker", "pull", "ubuntu")
}

func (p *BashProvider) Run(r *RunEnvironment) {
	cmd(append([]string{"/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "ubuntu", "bash", r.Cmd}, r.Args...)...)
}
