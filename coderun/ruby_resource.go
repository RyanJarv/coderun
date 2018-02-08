package coderun

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func RubyProvider() *Provider {
	return &Provider{
		RegisterOnCmd: rubyRegisterOnCmd,
		Setup:         rubySetup,
		Run:           rubyRun,
	}
}

func rubyRegisterOnCmd(cmd ...string) bool {
	match, err := regexp.MatchString(`^(([^ ]+/)?ruby[0-9.]* .*|[\S]+\.rb)$`, strings.Join(cmd, " "))
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func rubySetup(r *RunEnvironment) {
	r.CRDocker.Pull("ruby:2.3")

	if _, err := os.Stat("./Gemfile"); err == nil {
		Cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", r.CRDocker.newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "ruby:2.3", "bundler", "install", "--path", ".coderun/vendor/bundle")
	}
}

func rubyRun(r *RunEnvironment) {
	r.CRDocker.Run(dockerRunConfig{Image: "ruby:2.3", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{"ruby"}, r.Cmd...)})
}
