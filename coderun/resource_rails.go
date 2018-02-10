package coderun

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func RailsResource() *Resource {
	return &Resource{
		Register: railsRegister,
		Setup:    railsSetup,
		Run:      railsRun,
	}
}

func railsRegister(r RunEnvironment, p IProviderEnv) bool {
	match, err := regexp.MatchString(".*/?rails (.+ )?server.*", strings.Join(r.Cmd, " "))
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func railsSetup(r RunEnvironment, p IProviderEnv) {
	p.(dockerProviderEnv).CRDocker.Pull("ruby:2.3")

	image := p.(dockerProviderEnv).CRDocker.getOrBuildImage("ruby:2.3", []string{"sh", "-c", "apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y nodejs && bundler config --local path .coderun/vendor/bundle"})

	if _, err := os.Stat("./Gemfile"); err == nil {
		p.(dockerProviderEnv).CRDocker.Run(dockerRunConfig{Image: image, DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: []string{"bundler", "install", "--path", ".coderun/vendor/bundle"}})
	}
}

func railsRun(r RunEnvironment, p IProviderEnv) {
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
	image := p.(dockerProviderEnv).CRDocker.getImageName()
	log.Printf("Cmd is: %s", r.Cmd)
	p.(dockerProviderEnv).CRDocker.Run(dockerRunConfig{Image: image, DestDir: "/usr/src/myapp", SourceDir: Cwd(), Port: port, Cmd: append(append([]string{"bundler", "exec"}, r.Cmd...), "-b", "0.0.0.0")})
}
