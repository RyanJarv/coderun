package coderun

import (
	"fmt"
	"os"
)

func NewPythonProvider(conf *ProviderConfig) (Provider, error) {
	return &PythonProvider{
		name: "python",
	}, nil
}

type PythonProvider struct {
	name string
}

func (p *PythonProvider) Setup(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "pull", "python:3")

	if _, err := os.Stat("./requirements.txt"); err == nil {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "python:3", "python", "-m", "venv", ".coderun/venv")
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "python:3", "sh", "-c", ". .coderun/venv/bin/activate && pip install -r ./requirements.txt")
	}
}

func (p *PythonProvider) Run(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "python:3", "sh", "-c", fmt.Sprintf(". .coderun/venv/bin/activate && python %s", r.FilePath))
}
