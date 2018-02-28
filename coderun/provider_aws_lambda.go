package coderun

import (
	"log"

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
			"pipResource":     PipResource(),
			"awsLambdaPython": AWSLambdaPython(),
			"awsLambdaJs":     AWSLambdaJs(),
		},
		ProviderEnv: &awsLambdaProviderEnv{},
	}
}
func awsLambdaRegister(r *RunEnvironment) bool {
	return false
}

func awsLambdaResourceRegister(p Provider, runEnv *RunEnvironment) {
	for n, r := range p.Resources {
		if r.Register(runEnv, p) {
			p.RegisteredResources[n] = r
		}
	}

	if len(p.RegisteredResources) < 1 {
		log.Fatalf("Didn't find any registered lambda resources")
	}
}

func awsLambdaSetup(provider Provider, r *RunEnvironment) IProviderEnv {
	awsConfig := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	crLambda := NewCRLambda(awsConfig)
	providerEnv := awsLambdaProviderEnv{
		CRLambda: crLambda,
	}

	for _, resource := range provider.RegisteredResources {
		if resource.Setup == nil {
			Logger.info.Printf("No step Setup found for resource %s", resource.Name)
		} else {
			resource.Setup(r, providerEnv)
		}
	}

	return providerEnv
}

func awsLambdaDeploy(provider Provider, r *RunEnvironment, p IProviderEnv) {
	providerEnv := p.(awsLambdaProviderEnv)

	for _, resource := range provider.RegisteredResources {
		if resource.Deploy == nil {
			Logger.info.Printf("No step Deploy found for resource %s", resource.Name)
		} else {
			Logger.info.Printf("Running lambda resource %s", provider.Name)
			resource.Deploy(r, providerEnv)
		}
	}
}

func awsLambdaRun(provider Provider, r *RunEnvironment, p IProviderEnv) {
	providerEnv := p.(awsLambdaProviderEnv)

	for _, resource := range provider.RegisteredResources {
		if resource.Run == nil {
			Logger.info.Printf("No step Run found for resource %s", resource.Name)
		} else {
			Logger.info.Printf("Running lambda resource %s", provider.Name)
			resource.Run(r, providerEnv)
		}
	}
}
