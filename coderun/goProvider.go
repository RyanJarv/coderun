package coderun

import (
	"fmt"
	"log"
	"os"
)

func NewGoProvider(conf *ProviderConfig) (Provider, error) {
	return &GoProvider{
		name: "ruby",
	}, nil
}

type GoProvider struct {
	name string
}

func (p *GoProvider) Setup(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "pull", "golang")

	var image string
	if image, _ = getImageName(); image == "" {
		var image = newImageName()
		var preimage = newImageName()
		cmd("/usr/local/bin/docker", "run", "-t", "--name", preimage, "-v", fmt.Sprintf("%s:/go/src/localhost/myapp", r.Cwd), "-w", "/go/src/localhost/myapp", "golang", "sh", "-c", "curl https://glide.sh/get | sh")
		cmd("/usr/local/bin/docker", "commit", preimage, image)
		cmd("/usr/local/bin/docker", "rm", preimage)
		err := setImageName(image)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("./glide.lock"); os.IsNotExist(err) {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/go/src/localhost/myapp", r.Cwd), "-w", "/go/src/localhost/myapp", image, "glide", "init", "--non-interactive")
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/go/src/localhost/myapp", r.Cwd), "-w", "/go/src/localhost/myapp", image, "glide", "install")
	}
}

func (p *GoProvider) Run(r *RunEnvironment) {
	cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/go/src/localhost/myapp", r.Cwd), "-w", "/go/src/localhost/myapp", "golang", "go", "run", r.FilePath)
}
