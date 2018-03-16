package coderun

type ISnitchResource interface {
	Name() string
	Register(IRunEnvironment, IProvider) bool
	Setup(*StepCallback, *StepCallback)
}

func NewSnitchProvider(e IRunEnvironment) IProvider {
	return &SnitchProvider{
		resources: []ISnitchResource{
			NewSnitchDockerResource(e),
		},
	}
}

type SnitchProvider struct {
	resources []ISnitchResource
}

func (p *SnitchProvider) Name() string {
	return "snitch"
}

func (p *SnitchProvider) Register(e IRunEnvironment) bool {
	registered := false
	for _, r := range p.resources {
		if r.Register(e, p) == true {
			registered = true
		}
	}
	return registered
}
