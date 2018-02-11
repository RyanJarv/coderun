package coderun

func AWSLambdaPython() *Resource {
	return &Resource{
		Register: awsLambdaPythonRegister,
		Setup:    awsLambdaPythonSetup,
		Run:      awsLambdaPythonRun,
	}
}

func awsLambdaPythonRegister(r RunEnvironment, p IProviderEnv) bool {
	return true
}

func awsLambdaPythonSetup(r RunEnvironment, p IProviderEnv) {
	//p.(awsLambdaProviderEnv).CRDocker.Pull("bash")
}

func awsLambdaPythonRun(r RunEnvironment, p IProviderEnv) {
	//p.(awsLambdaProviderEnv).CRDocker.Run(dockerRunConfig{Image: "ubuntu", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{"bash"}, r.Cmd...)})
}
