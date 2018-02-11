package coderun

import (
	"fmt"
	"log"
	"os"
)

func PythonResource() *Resource {
	return &Resource{
		Register: pythonRegister,
		Setup:    pythonSetup,
		Run:      pythonRun,
	}
}

func pythonRegister(r RunEnvironment, p IProviderEnv) bool {
	log.Printf("Command is %s", r.Cmd)
	return MatchCommandOrExt(r.Cmd, "python", ".py")
}

func pythonSetup(r RunEnvironment, p IProviderEnv) {
	p.(dockerProviderEnv).CRDocker.Pull("python:3")

	Exec("/usr/local/bin/docker", "run", "-t", "--rm", "--name", p.(dockerProviderEnv).CRDocker.newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "python:3", "python", "-m", "venv", ".coderun/venv")

	if _, err := os.Stat("./requirements.txt"); err == nil {
		Exec("/usr/local/bin/docker", "run", "-t", "--rm", "--name", p.(dockerProviderEnv).CRDocker.newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "python:3", "sh", "-c", ". .coderun/venv/bin/activate && pip install -r ./requirements.txt")
	}
}

func pythonRun(r RunEnvironment, p IProviderEnv) {
	log.Printf("Args: %s", r.Cmd)
	p.(dockerProviderEnv).CRDocker.Run(dockerRunConfig{Image: "python:3", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{".coderun/venv/bin/python"}, r.Cmd...)})
}
