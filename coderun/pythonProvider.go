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
	dockerPull(r.DockerClient, "python:3")

	if _, err := os.Stat("./requirements.txt"); err == nil {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "python:3", "python", "-m", "venv", ".coderun/venv")
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "python:3", "sh", "-c", ". .coderun/venv/bin/activate && pip install -r ./requirements.txt")
	}
}

func (p *PythonProvider) Run(r *RunEnvironment) {
	log.Printf("Args: %s", r.Cmd)
	dockerRun(dockerRunConfig{Client: r.DockerClient, Image: "python:3", DestDir: "/usr/src/myapp", SourceDir: r.Cwd, Cmd: append([]string{".coderun/venv/bin/python"}, r.Cmd...)})
}
