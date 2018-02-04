package coderun

import (
	"log"
	"os"
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
	dockerPull("ruby:2.3")

	image := getOrBuildImage("ruby:2.3", []string{"sh", "-c", "apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y nodejs && bundler config --local path .coderun/vendor/bundle"})

	if _, err := os.Stat("./Gemfile"); err == nil {
		dockerRun(image, 1234, "/go/src/localhost/myapp", "bundler", "install", "--path", ".coderun/vendor/bundle")
	}
}

func (p *RailsProvider) Run(r *RunEnvironment) {
	port := 3000
	var err error
	for i, arg := range r.Args {
		if arg == `-p` {
			port, err = strconv.Atoi(r.Args[i+1])
			if err != nil {
				log.Fatal(err)
			}
		}
		if arg == `-b` {
			if r.Args[i+1] != "0.0.0.0" {
				log.Fatal("Rails server bind IP needs to be 0.0.0.0 to be reachable from localhost (note this does not actually expose it publicly when running in docker)")
			}
		}
	}
	image := getImageName()
	log.Printf("Args is: %s", r.Args)
	dockerRun(image, port, "/usr/local/myapp", append([]string{"bundler", "exec", r.Cmd}, append(r.Args, []string{"-b", "0.0.0.0"}...)...)...)
}
