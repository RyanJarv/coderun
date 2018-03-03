package coderun

func NewBashResource() IDockerResource {
	return &BashResource{}
}

type BashResource struct {
	IResource
}

func (r *BashResource) Name() string {
	return "ubuntu"
}

func (r *BashResource) Register(e *RunEnvironment, p IProvider) bool {
	if MatchCommandOrExt(e.Cmd, "bash", ".sh") {
		e.Registry.AddAt(SetupStep, &StepCallback{Step: "Setup", Provider: p, Resource: r, Callback: r.Setup})
		e.Registry.AddAt(SetupStep, &StepCallback{Step: "Run", Provider: p, Resource: r, Callback: r.Run})
		return true
	} else {
		return false
	}
}

func (r *BashResource) RegisterMount(e *RunEnvironment, local string, remote string) {
	e.CRDocker.RegisterMount(local, remote)
}

func (r *BashResource) Setup(e *RunEnvironment, callback *StepCallback, currentStep *StepCallback) {
	e.CRDocker.Pull("ubuntu")
}

func (r *BashResource) Run(e *RunEnvironment, callback *StepCallback, currentStep *StepCallback) {
	e.CRDocker.Run(dockerRunConfig{
		Image:     "ubuntu",
		DestDir:   "/usr/src/myapp",
		SourceDir: Cwd(),
		Cmd:       append([]string{"bash"}, e.Cmd...),
	})
}
