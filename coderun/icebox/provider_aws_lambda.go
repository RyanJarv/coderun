package coderun

import (
	"github.com/RyanJarv/coderun/coderun"
	"log"

	"github.com/aws/aws-sdk-go/aws"
)

type AwsLambdaProviderEnv struct {
	IProviderEnv
	CRLambda *coderun.CRLambda
}

func AWSLambdaProvider() *Provider {
	return &Provider{
		Name:                "awsLambda",
		Register:            awsLambdaRegister,
		ResourceRegister:    awsLambdaResourceRegister,
		Setup:               awsLambdaSetup,
		Deploy:              awsLambdaDeploy,
		Run:                 awsLambdaRun,
		RegisteredResources: map[string]IResource{},
		Resources: map[string]IResource{
			"awsLambdaPython": AWSLambdaPython(),
		},
		ProviderEnv: &AwsLambdaProviderEnv{},
	}
}
func awsLambdaRegister(r *RunEnvironment) bool {
	return false
}

func awsLambdaResourceRegister(p *Provider, runEnv *RunEnvironment) {
	for name, r := range p.Resources {
		if r.(Resource).Register(runEnv, p) {
			Logger.info.Printf("Registering resource %s", name)
			p.RegisteredResources[name] = r
		}
	}

	if len(p.RegisteredResources) < 1 {
		log.Fatalf("Didn't find any registered lambda resources")
	}
}

func awsLambdaSetup(p *Provider, r *RunEnvironment) IProviderEnv {
	awsConfig := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	crLambda := coderun.NewCRLambda(awsConfig)
	providerEnv := AwsLambdaProviderEnv{
		CRLambda: crLambda,
	}

	for _, resource := range p.RegisteredResources {
		if resource.(Resource).Setup == nil {
			Logger.info.Printf("No step Setup found for resource %s", resource.(Resource).Name)
		} else {
			resource.(Resource).Setup(r, providerEnv)
		}
	}

	return providerEnv
}

func awsLambdaDeploy(provider *Provider, r *RunEnvironment, p IProviderEnv) {
	providerEnv := p.(AwsLambdaProviderEnv)

	for _, resource := range provider.RegisteredResources {
		if resource.(Resource).Deploy == nil {
			Logger.info.Printf("No step Deploy found for resource %s", resource.(Resource).Name)
		} else {
			Logger.info.Printf("Running lambda resource %s", provider.Name)
			resource.(Resource).Deploy(r, providerEnv)
		}
	}
}

func awsLambdaRun(provider *Provider, r *RunEnvironment, p IProviderEnv) {
	providerEnv := p.(AwsLambdaProviderEnv)

	for _, resource := range provider.RegisteredResources {
		if resource.(Resource).Run == nil {
			Logger.info.Printf("No step Run found for resource %s", resource.(Resource).Name)
		} else {
			Logger.info.Printf("Running lambda resource %s", provider.Name)
			resource.(Resource).Run(r, providerEnv)
		}
	}
}
