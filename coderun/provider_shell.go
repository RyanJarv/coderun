package coderun

type IShellResource interface {
	IResource
	Setup(*StepCallback, *StepCallback)
	Run(*StepCallback, *StepCallback)
}

type ShellResources map[string]IResource

func NewShellProvider(r IRunEnvironment) IProvider {
	return &ShellProvider{
		resources: map[string]IShellResource{
			"bash": NewShellBashResource(r),
		},
		registeredResources: map[string]IShellResource{},
	}
}

type ShellProvider struct {
	IProvider
	resources           map[string]IShellResource
	registeredResources map[string]IShellResource
}

func (p *ShellProvider) Name() string {
	return "shell"
}

func (p *ShellProvider) Register(e IRunEnvironment) bool {
	registered := false
	for name, r := range p.resources {
		if r.Register(e, p) {
			Logger.info.Printf("Registering resource %s", name)
			p.registeredResources[name] = r
			registered = true
		}
	}
	return registered
}

func (p *ShellProvider) Setup(callback *StepCallback, currentStep *StepCallback) {
}

func (p *ShellProvider) Run(callback *StepCallback, currentStep *StepCallback) {
}

func (p *ShellProvider) Teardown(callback *StepCallback, currentStep *StepCallback) {
}
