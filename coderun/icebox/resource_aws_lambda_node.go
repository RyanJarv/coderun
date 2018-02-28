package coderun

func AWSLambdaJs() *Resource {
	return &Resource{
		Register: awsLambdaJsRegister,
		Setup:    awsLambdaJsSetup,
		Deploy:   awsLambdaJsDeploy,
		Run:      awsLambdaJsRun,
	}
}

func awsLambdaJsRegister(r *RunEnvironment, p IProviderEnv) bool {
	return MatchCommandOrExt(r.Cmd, "node", ".js")
}

func awsLambdaJsSetup(r *RunEnvironment, p IProviderEnv) {
	Logger.debug.Printf("r.DependsDir is %s", r.DependsDir)
	p.(awsLambdaProviderEnv).CRLambda.Setup(r)
}

func awsLambdaJsDeploy(r *RunEnvironment, p IProviderEnv) {
	pEnv := p.(awsLambdaProviderEnv)
	pEnv.CRLambda.Deploy("nodejs6.10", r, pEnv)
}

func awsLambdaJsRun(r *RunEnvironment, p IProviderEnv) {
	providerEnv := p.(awsLambdaProviderEnv)
	providerEnv.CRLambda.Run(r, providerEnv)
}
