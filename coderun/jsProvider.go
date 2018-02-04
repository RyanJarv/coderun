package coderun

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type JsProvider struct {
	ProviderDefault
	name string
}

func (p *JsProvider) RegisterOnCmd(cmd string, args ...string) bool {
	match, err := regexp.MatchString(`^(([^ ]+/)?node(js)? .*|[\S]+\.js)$`, strings.Join(append([]string{cmd}, args...), " "))
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func (p *JsProvider) Setup(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "pull", "node")

	if _, err := os.Stat("./package-lock.json"); os.IsNotExist(err) {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "node", "npm", "install")
	}
}

func (p *JsProvider) Run(r *RunEnvironment) {
	cmd(append([]string{"/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "node", "node", r.Cmd}, r.Args...)...)
}
