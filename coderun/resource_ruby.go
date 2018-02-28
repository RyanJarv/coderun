package coderun

import (
	"fmt"
	"os"
)

func RubyResource() *Resource {
	return &Resource{
		Register: rubyRegister,
		Setup:    rubySetup,
		Run:      rubyRun,
	}
}

func rubyRegister(r *RunEnvironment, p IProviderEnv) bool {
	return MatchCommandOrExt(r.Cmd, "ruby", ".rb")
}

func rubySetup(r *RunEnvironment, p IProviderEnv) {
	r.CRDocker.Pull("ruby:2.3")

	if _, err := os.Stat("./Gemfile"); err == nil {
		Exec("/usr/local/bin/docker", "run", "-t", "--rm", "--name", r.CRDocker.newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "ruby:2.3", "bundler", "install", "--path", ".coderun/vendor/bundle")
	}
}

func rubyRun(r *RunEnvironment, p IProviderEnv) {
	r.CRDocker.Run(dockerRunConfig{Image: "ruby:2.3", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{"ruby"}, r.Cmd...)})
}
