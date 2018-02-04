package coderun

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type PythonProvider struct {
	ProviderDefault
	name string
}

func (p *PythonProvider) RegisterOnCmd(cmd string, args ...string) bool {
	match, err := regexp.MatchString(`^(([^ ]+/)?python[0-9.]* .*|[\S]+\.py)$`, strings.Join(append([]string{cmd}, args...), " "))
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func (p *PythonProvider) Setup(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "pull", "python:3")

	if _, err := os.Stat("./requirements.txt"); err == nil {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "python:3", "python", "-m", "venv", ".coderun/venv")
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "python:3", "sh", "-c", ". .coderun/venv/bin/activate && pip install -r ./requirements.txt")
	}
}

func (p *PythonProvider) Run(r *RunEnvironment) {
	log.Printf("Args is: %s", r.Cmd)
	dockerRun("python:3", 1234, "/usr/local/myapp", append([]string{".coderun/venv/bin/python", r.Cmd}, r.Args...)...)
}
