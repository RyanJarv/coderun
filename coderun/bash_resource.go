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

func bashSetup(r RunEnvironment) {
	r.CRDocker.Pull("bash")
}

func bashRun(r RunEnvironment) {
	r.CRDocker.Run(dockerRunConfig{Image: "ubuntu", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{"bash"}, r.Cmd...)})
}
