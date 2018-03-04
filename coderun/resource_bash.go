package coderun

func NewBashResource(r **RunEnvironment) IDockerResource {
	return &BashResource{env: r}
}

type BashResource struct {
	IResource
	env **RunEnvironment
}

func (r *BashResource) Name() string {
	return "ubuntu"
}

func (r *BashResource) Register(p IProvider) bool {
	if MatchCommandOrExt((*r.env).Cmd, "bash", ".sh") {
		(*r.env).Registry.AddAt(SetupStep, &StepCallback{Step: "Setup", Provider: p, Resource: r, Callback: r.Setup})
		(*r.env).Registry.AddAt(SetupStep, &StepCallback{Step: "Run", Provider: p, Resource: r, Callback: r.Run})
		return true
	} else {
		return false
	}
}

func (r *BashResource) RegisterMount(local string, remote string) {
	(*r.env).CRDocker.RegisterMount(local, remote)
}

func (r *BashResource) Setup(callback *StepCallback, currentStep *StepCallback) {
	(*r.env).CRDocker.Pull("ubuntu")
}

func (r *BashResource) Run(callback *StepCallback, currentStep *StepCallback) {
	(*r.env).CRDocker.Run(dockerRunConfig{
		Image:     "ubuntu",
		DestDir:   "/usr/src/myapp",
		SourceDir: Cwd(),
		Cmd:       append([]string{"bash"}, (*r.env).Cmd...),
	})
}
