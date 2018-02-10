package coderun

func BashResource() *Resource {
	return &Resource{
		Register: bashRegister,
		Setup:    bashSetup,
		Run:      bashRun,
	}
}

func bashRegister(r RunEnvironment, p IProviderEnv) bool {
	return MatchCommandOrExt(r.Cmd, "bash", ".sh")
}

func bashSetup(r RunEnvironment, p IProviderEnv) {
	p.(dockerProviderEnv).CRDocker.Pull("bash")
}

func bashRun(r RunEnvironment, p IProviderEnv) {
	p.(dockerProviderEnv).CRDocker.Run(dockerRunConfig{Image: "ubuntu", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{"bash"}, r.Cmd...)})
}
