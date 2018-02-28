package coderun

func PythonResource() *Resource {
	return &Resource{
		Register: pythonRegister,
		Setup:    pythonSetup,
		Run:      pythonRun,
	}
}

func pythonRegister(r *RunEnvironment, p IProviderEnv) bool {
	return MatchCommandOrExt(r.Cmd, "python", ".py")
}

func pythonSetup(r *RunEnvironment, p IProviderEnv) {
	r.CRDocker.Pull("python:3")

	r.CRDocker.Run(dockerRunConfig{
		Image:     "python:3",
		DestDir:   "/usr/src/myapp",
		SourceDir: Cwd(),
		Cmd:       []string{"python", "-m", "venv", ".coderun/venv"},
	})
}

func pythonRun(r *RunEnvironment, p IProviderEnv) {
	r.CRDocker.Run(dockerRunConfig{
		Image:     "python:3",
		DestDir:   "/usr/src/myapp",
		SourceDir: Cwd(),
		Cmd:       append([]string{".coderun/venv/bin/python"}, r.Cmd...),
	})
}
