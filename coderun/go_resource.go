package coderun

import (
	"os"
)

func GoResource() *Resource {
	return &Resource{
		Register: goRegister,
		Setup:    goSetup,
		Run:      goRun,
	}
}

func goRegister(r *RunEnvironment, p IProviderEnv) bool {
	return MatchCommandOrExt(r.Cmd, "go", ".go")
}

func goSetup(r *RunEnvironment, d IProviderEnv) {
	r.CRDocker.Pull("golang")

	if _, err := os.Stat("./glide.lock"); os.IsNotExist(err) {
		r.CRDocker.Run(dockerRunConfig{Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: Cwd(), Cmd: []string{"sh", "-c", "curl https://glide.sh/get | sh && glide init --non-interactive && glide install"}})
	}
}

func goRun(r *RunEnvironment, d IProviderEnv) {
	r.CRDocker.Run(dockerRunConfig{Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: Cwd(), Cmd: []string{"go", "run"}})
}
