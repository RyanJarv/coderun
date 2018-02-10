package coderun

import (
	"fmt"
	"log"
	"os"
)

func PythonResource() *Resource {
	return &Resource{
		RegisterOnCmd: pythonRegisterOnCmd,
		Setup:         pythonSetup,
		Run:           pythonRun,
	}
}

func pythonRegisterOnCmd(cmd ...string) bool {
	return MatchCommandOrExt(cmd, "python", ".py")
}

func pythonSetup(r RunEnvironment) {
	r.CRDocker.Pull("python:3")

	if _, err := os.Stat("./requirements.txt"); err == nil {
		r.Exec("/usr/local/bin/docker", "run", "-t", "--rm", "--name", r.CRDocker.(CRDocker).newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "python:3", "python", "-m", "venv", ".coderun/venv")
		r.Exec("/usr/local/bin/docker", "run", "-t", "--rm", "--name", r.CRDocker.(CRDocker).newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "python:3", "sh", "-c", ". .coderun/venv/bin/activate && pip install -r ./requirements.txt")
	}
}

func pythonRun(r RunEnvironment) {
	log.Printf("Args: %s", r.Cmd)
	r.CRDocker.Run(dockerRunConfig{Image: "python:3", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{".coderun/venv/bin/python"}, r.Cmd...)})
}
