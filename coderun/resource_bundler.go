package coderun

import (
	"fmt"
	"os"
	"path/filepath"
)

func BundlerResource() *Resource {
	return &Resource{
		Register: bundlerRegister,
		Setup:    bundlerSetup,
	}
}

func bundlerRegister(r *RunEnvironment, p IProviderEnv) bool {
	f := filepath.Join(r.CodeDir, "Gemfile")
	Logger.debug.Printf("Checking for file at %s", f)
	_, err := os.Stat(f)
	return MatchCommandOrExt(r.Cmd, "ruby", ".rb") && err == nil
}

func bundlerSetup(r *RunEnvironment, p IProviderEnv) {
	r.CRDocker.Pull("ruby:2.3")
	if _, err := os.Stat("./Gemfile"); err == nil {
		Exec("/usr/local/bin/docker", "run", "-t", "--rm", "--name", r.CRDocker.newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "ruby:2.3", "bundle", "install", "--path", ".coderun/vendor/bundle")
	}
}
