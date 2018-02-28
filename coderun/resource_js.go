package coderun

import (
	"fmt"
	"os"
)

func JsResource() *Resource {
	return &Resource{
		Register: jsRegister,
		Setup:    jsSetup,
		Run:      jsRun,
	}
}

func jsRegister(r *RunEnvironment, p IProviderEnv) bool {
	return MatchCommandOrExt(r.Cmd, "node", ".js")
}

func jsSetup(r *RunEnvironment, p IProviderEnv) {
	r.Exec("/usr/local/bin/docker", "pull", "node")

	if _, err := os.Stat("./package-lock.json"); os.IsNotExist(err) {
		r.Exec("/usr/local/bin/docker", "run", "-t", "--rm", "--name", r.CRDocker.newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "node", "npm", "install")
	}
}

func jsRun(r *RunEnvironment, p IProviderEnv) {
	r.CRDocker.Run(dockerRunConfig{Image: "node", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{"node"}, r.Cmd...)})
}
