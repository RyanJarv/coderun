package coderun

func AWSLambdaPython() *Resource {
	return &Resource{
		Register: awsLambdaPythonRegister,
		Setup:    awsLambdaPythonSetup,
		Deploy:   awsLambdaPythonDeploy,
		Run:      awsLambdaPythonRun,
	}
}

func awsLambdaPythonRegister(r *RunEnvironment, p IProviderEnv) bool {
	return true
}

func awsLambdaPythonSetup(r *RunEnvironment, p IProviderEnv) {
	Logger.debug.Printf("r.DependsDir is %s", r.DependsDir)
	if r.DependsDir == "" {
		panic("asdf")
	}
	p.(awsLambdaProviderEnv).CRLambda.Setup(r)
}

func awsLambdaPythonDeploy(r *RunEnvironment, p IProviderEnv) {
	pEnv := p.(awsLambdaProviderEnv)
	pEnv.CRLambda.Deploy(r, pEnv)
}

func awsLambdaPythonRun(r *RunEnvironment, p IProviderEnv) {
	//p.(awsLambdaProviderEnv).CRDocker.Run(dockerRunConfig{Image: "ubuntu", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{"bash"}, r.Cmd...)})
}
