package coderun

import (
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
	return MatchCommandOrExt(r.Cmd, "python", ".py")
}

func pythonSetup(r RunEnvironment, p IProviderEnv) {
	p.(dockerProviderEnv).CRDocker.Pull("python:3")

	p.(dockerProviderEnv).CRDocker.Run(dockerRunConfig{
		Image:     "python:3",
		DestDir:   "/usr/src/myapp",
		SourceDir: Cwd(),
		Cmd:       []string{"python", "-m", "venv", ".coderun/venv"},
	})

	if _, err := os.Stat("./requirements.txt"); err == nil {
		p.(dockerProviderEnv).CRDocker.Run(dockerRunConfig{
			Image:     "python:3",
			DestDir:   "/usr/src/myapp",
			SourceDir: Cwd(),
			Cmd:       []string{"sh", "-c", ". .coderun/venv/bin/activate && pip install -r ./requirements.txt"},
		})
	}
}

func pythonRun(r RunEnvironment, p IProviderEnv) {
	p.(dockerProviderEnv).CRDocker.Run(dockerRunConfig{
		Image:     "python:3",
		DestDir:   "/usr/src/myapp",
		SourceDir: Cwd(),
		Cmd:       append([]string{".coderun/venv/bin/python"}, r.Cmd...),
	})
}
