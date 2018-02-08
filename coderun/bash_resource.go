package coderun

func BashResource() *Resource {
	return &Resource{
		RegisterOnCmd: bashRegisterOnCmd,
		Setup:         bashSetup,
		Run:           bashRun,
	}
}

func bashRegisterOnCmd(cmd ...string) bool {
	return MatchCommandOrExt(cmd, "bash", ".sh")
}

func bashSetup(r IRunEnvironment) {
	r.(RunEnvironment).CRDocker.Pull("bash")
}

func bashRun(r IRunEnvironment) {
	r.(RunEnvironment).CRDocker.Run(dockerRunConfig{Image: "ubuntu", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{"bash"}, r.Cmd()...)})
}
