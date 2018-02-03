package coderun

import (
	"fmt"
	"os"
)

func NewJsProvider(conf *ProviderConfig) (Provider, error) {
	return &JsProvider{
		name: "ruby",
	}, nil
}

type JsProvider struct {
	name string
}

func (p *JsProvider) Setup(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "pull", "node")

	if _, err := os.Stat("./package-lock.json"); os.IsNotExist(err) {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "node", "npm", "install")
	}
}

func (p *JsProvider) Run(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", r.Cwd), "-w", "/usr/src/myapp", "node", "node", r.FilePath)
}
