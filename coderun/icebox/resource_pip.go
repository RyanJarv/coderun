package coderun

import (
	"os"
	"path/filepath"
)

func PipResource() *Resource {
	return &Resource{
		Register: pipRegister,
		Setup:    pipSetup,
	}
}

func pipRegister(r *RunEnvironment, p IProviderEnv) bool {
	f := filepath.Join(r.CodeDir, "requirements.txt")
	Logger.debug.Printf("Checking for file at %s", f)
	_, err := os.Stat(f)
	return MatchCommandOrExt(r.Cmd, "python", ".py") && err == nil
}

func pipSetup(r *RunEnvironment, p IProviderEnv) {
	Logger.debug.Printf("Running pipSetup")
	r.CRDocker.Pull("python:3")
	r.DependsDir = ".coderun/venv/lib/python3.6/site-packages/"
	r.CRDocker.Run(dockerRunConfig{
		Image:     "python:3",
		DestDir:   "/usr/src/myapp",
		SourceDir: Cwd(),
		Cmd:       []string{"python", "-m", "venv", ".coderun/venv"},
	})

	if _, err := os.Stat("./requirements.txt"); err == nil {
		r.CRDocker.Run(dockerRunConfig{
			Image:     "python:3",
			DestDir:   "/usr/src/myapp",
			SourceDir: Cwd(),
			Cmd:       []string{"sh", "-c", ". .coderun/venv/bin/activate && pip install -r ./requirements.txt"},
		})
	}
}
