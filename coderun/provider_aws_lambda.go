package coderun

import (
	"github.com/aws/aws-sdk-go/aws"
)

type awsLambdaProviderEnv struct {
	IProviderEnv
	CRLambda *CRLambda
}

func AWSLambdaProvider() *Provider {
	return &Provider{
		Name:                "awsLambda",
		Register:            awsLambdaRegister,
		ResourceRegister:    awsLambdaResourceRegister,
		Setup:               awsLambdaSetup,
		Deploy:              awsLambdaDeploy,
		Run:                 awsLambdaRun,
		RegisteredResources: map[string]*Resource{},
		Resources: map[string]*Resource{
			"awsLambdaPython": AWSLambdaPython(),
		},
		ProviderEnv: &awsLambdaProviderEnv{},
	}
}
func awsLambdaRegister(r RunEnvironment) bool {
	return false
}

func awsLambdaResourceRegister(p Provider, runEnv RunEnvironment) {
	for n, r := range p.Resources {
		if r.Register(runEnv, p) {
			p.RegisteredResources[n] = r
		}
	}
}

func awsLambdaSetup(provider Provider, r RunEnvironment) IProviderEnv {
	awsConfig := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	crLambda := NewCRLambda(awsConfig)
	crLambda.Setup(r)

	providerEnv := awsLambdaProviderEnv{
		CRLambda: crLambda,
	}
	return providerEnv
}

func awsLambdaDeploy(provider Provider, r RunEnvironment, p IProviderEnv) {
	providerEnv := p.(awsLambdaProviderEnv)
	providerEnv.CRLambda.Deploy(r, providerEnv)
}

func awsLambdaRun(provider Provider, r RunEnvironment, p IProviderEnv) {
	providerEnv := p.(awsLambdaProviderEnv)
	providerEnv.CRLambda.Run(r, providerEnv)
}
