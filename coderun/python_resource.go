package coderun

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func PythonProvider() *Provider {
	return &Provider{
		RegisterOnCmd: pythonRegisterOnCmd,
		Setup:         pythonSetup,
		Run:           pythonRun,
	}
}

func pythonRegisterOnCmd(cmd ...string) bool {
	match, err := regexp.MatchString(`^(([^ ]+/)?python[0-9.]* .*|[\S]+\.py)$`, strings.Join(cmd, " "))
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func pythonSetup(r *RunEnvironment) {
	r.CRDocker.Pull("python:3")

	if _, err := os.Stat("./requirements.txt"); err == nil {
		Cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", r.CRDocker.newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "python:3", "python", "-m", "venv", ".coderun/venv")
		Cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", r.CRDocker.newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "python:3", "sh", "-c", ". .coderun/venv/bin/activate && pip install -r ./requirements.txt")
	}
}

func pythonRun(r *RunEnvironment) {
	log.Printf("Args: %s", r.Cmd)
	r.CRDocker.Run(dockerRunConfig{Image: "python:3", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{".coderun/venv/bin/python"}, r.Cmd...)})
}
