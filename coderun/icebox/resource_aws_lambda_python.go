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
	return MatchCommandOrExt(r.Cmd, "python", ".py")
}

func awsLambdaPythonSetup(r *RunEnvironment, p IProviderEnv) {
	Logger.debug.Printf("r.DependsDir is %s", r.DependsDir)
	p.(awsLambdaProviderEnv).CRLambda.Setup(r)
}

func awsLambdaPythonDeploy(r *RunEnvironment, p IProviderEnv) {
	pEnv := p.(awsLambdaProviderEnv)
	pEnv.CRLambda.Deploy("python3.6", r, pEnv)
}

func awsLambdaPythonRun(r *RunEnvironment, p IProviderEnv) {
	providerEnv := p.(awsLambdaProviderEnv)
	providerEnv.CRLambda.Run(r, providerEnv)
}
