package coderun

func RubyResource() *Resource {
	return &Resource{
		Register: rubyRegister,
		Setup:    rubySetup,
		Run:      rubyRun,
	}
}

func rubyRegister(r *RunEnvironment, p IProviderEnv) bool {
	Logger.debug.Printf("Running rubyRegister")
	return MatchCommandOrExt(r.Cmd, "ruby", ".rb")
}

func rubySetup(r *RunEnvironment, p IProviderEnv) {
	r.CRDocker.Pull("ruby:2.3")
}

func rubyRun(r *RunEnvironment, p IProviderEnv) {
	r.CRDocker.Run(dockerRunConfig{Image: "ruby:2.3", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{"ruby"}, r.Cmd...)})
}
