package coderun

import (
	"github.com/aws/aws-sdk-go/aws"
	"log"
)

type ILambdaResource interface {
	IResource
}

type LambdaProvider struct {
	Provider
}

func NewLambdaProvider(r IRunEnvironment) IProvider {
	return &LambdaProvider{
		Provider{
			name: "lambda",
			resources: map[string]ILambdaResource{
				"bash": NewBashResource(r),
			},
		},
	}
}

func (p *LambdaProvider) Register(e IRunEnvironment) bool {
	registered := false
	e.Registry().AddAt(TeardownStep+10, &StepCallback{Step: "Teardown", Provider: p, Callback: p.Teardown})
	for name, r := range p.resources {
		if r.Register(e, p) {
			Logger.Info.Printf("Registering resource %s", name)
			p.registeredResources[name] = r
			registered = true
		}
	}
	if registered == true {
		// This can be removed when buildkit get's merged into docker
		//(*p.env).Registry.AddAt(SetupStep-10, &StepCallback{Step: "Setup", Provider: p, Callback: p.Setup})
		//(*p.env).Registry.AddAt(TeardownStep, &StepCallback{Step: "Teardown", Provider: p, Callback: p.Teardown})
	}
	return registered
}

func (p *LambdaProvider) Setup(callback *StepCallback, currentStep *StepCallback) {
	awsConfig := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	crLambda := NewCRLambda(awsConfig)
	providerEnv := awsLambdaProviderEnv{
		CRLambda: crLambda,
	}

	// Should run with call backs the same way the provider resource does
	for _, resource := range p.registeredResources {
		if resource.(ILambdaResource).Setup == nil {
			Logger.Info.Printf("No step Setup found for resource %s", resource.(Resource).Name)
		} else {
			resource.(IResource).Setup(r, providerEnv)
		}
	}

	return providerEnv
}

func (p *LambdaProvider) Teardown(callback *StepCallback, currentStep *StepCallback) {
	NewCRLambda().LambdaKillLabel("coderun")
}
