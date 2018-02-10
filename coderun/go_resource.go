package coderun

import (
	"os"
)

func GoResource() *Resource {
	return &Resource{
		RegisterOnCmd: goRegisterOnCmd,
		Setup:         goSetup,
		Run:           goRun,
	}
}

func goRegisterOnCmd(cmd ...string) bool {
	return MatchCommandOrExt(cmd, "go", ".py")
}

func goSetup(r RunEnvironment) {
	r.CRDocker.(ICRDocker).Pull("golang")

	if _, err := os.Stat("./glide.lock"); os.IsNotExist(err) {
		r.CRDocker.Run(dockerRunConfig{Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: Cwd(), Cmd: []string{"sh", "-c", "curl https://glide.sh/get | sh && glide init --non-interactive && glide install"}})
	}
}

func goRun(r RunEnvironment) {
	r.CRDocker.Run(dockerRunConfig{Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: Cwd(), Cmd: []string{"go", "run"}})
}
