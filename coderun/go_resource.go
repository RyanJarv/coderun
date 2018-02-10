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

func goRegister(r RunEnvironment, p IProviderEnv) bool {
	return MatchCommandOrExt(r.Cmd, "go", ".py")
}

func goSetup(r RunEnvironment, d IProviderEnv) {
	d.(dockerProviderEnv).CRDocker.Pull("golang")

	if _, err := os.Stat("./glide.lock"); os.IsNotExist(err) {
		d.(dockerProviderEnv).CRDocker.Run(dockerRunConfig{Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: Cwd(), Cmd: []string{"sh", "-c", "curl https://glide.sh/get | sh && glide init --non-interactive && glide install"}})
	}
}

func goRun(r RunEnvironment, d IProviderEnv) {
	d.(dockerProviderEnv).CRDocker.Run(dockerRunConfig{Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: Cwd(), Cmd: []string{"go", "run"}})
}
