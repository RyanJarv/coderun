package coderun

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type RailsProvider struct {
	ProviderDefault
}

func (p *RailsProvider) RegisterOnCmd(cmd string, args ...string) bool {
	match, err := regexp.MatchString(".*/?rails (.+ )?server.*", strings.Join(append([]string{cmd}, args...), " "))
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func (p *RailsProvider) Setup(r *RunEnvironment) {
	//dockerPull(r.DockerClient, "ruby:2.3")

	//image := getOrBuildImage("ruby:2.3", []string{"sh", "-c", "apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y nodejs && bundler config --local path .coderun/vendor/bundle"})

	//if _, err := os.Stat("./Gemfile"); err == nil {
	//	dockerRun(dockerRunConfig{Client: r.DockerClient, Image: image, DestDir: "/usr/src/myapp", SourceDir: r.Cwd, Cmd: []string{"bundler", "install", "--path", ".coderun/vendor/bundle"}})
	//}
}

func (p *RailsProvider) Run(r *RunEnvironment) {
	var port int = 3000
	var err error
	for i, arg := range r.Cmd {
		if arg == `-p` {
			port, err = strconv.Atoi(r.Cmd[i+1])
			if err != nil {
				log.Fatal(err)
			}
		}
		if arg == `-b` {
			if r.Cmd[i+1] != "0.0.0.0" {
				log.Fatal("Rails server bind IP needs to be 0.0.0.0 to be reachable from localhost (note this does not actually expose it publicly when running in docker)")
			}
		}
	}
	log.Printf("%d", port)
	image := getImageName()
	log.Printf("Cmd is: %s", r.Cmd)
	dockerRun(dockerRunConfig{Client: r.DockerClient, Image: image, DestDir: "/usr/src/myapp", SourceDir: r.Cwd, Port: port, Cmd: append(append([]string{"bundler", "exec"}, r.Cmd...), "-b", "0.0.0.0")})
}
