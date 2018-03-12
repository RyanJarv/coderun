package coderun

type ISnitchResource interface {
	Name() string
	Register(IProvider) bool
	Setup(*StepCallback, *StepCallback)
}

func NewSnitchProvider(r **RunEnvironment) IProvider {
	return &SnitchProvider{
		resources: []ISnitchResource{
			NewSnitchDockerResource(r),
		},
	}
}

type SnitchProvider struct {
	resources []ISnitchResource
}

func (p *SnitchProvider) Name() string {
	return "snitch"
}

func (p *SnitchProvider) Register() bool {
	registered := false
	for _, r := range p.resources {
		if r.Register(p) == true {
			registered = true
		}
	}
	return registered
}
