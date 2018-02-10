package coderun

import (
	"fmt"
	"os"
)

func RubyResource() *Resource {
	return &Resource{
		RegisterOnCmd: rubyRegisterOnCmd,
		Setup:         rubySetup,
		Run:           rubyRun,
	}
}

func rubyRegisterOnCmd(cmd ...string) bool {
	return MatchCommandOrExt(cmd, "ruby", ".rb")
}

func rubySetup(r RunEnvironment) {
	r.CRDocker.Pull("ruby:2.3")

	if _, err := os.Stat("./Gemfile"); err == nil {
		r.Exec("/usr/local/bin/docker", "run", "-t", "--rm", "--name", r.CRDocker.newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "ruby:2.3", "bundler", "install", "--path", ".coderun/vendor/bundle")
	}
}

func rubyRun(r RunEnvironment) {
	r.CRDocker.Run(dockerRunConfig{Image: "ruby:2.3", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{"ruby"}, r.Cmd...)})
}
