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
		e.Registry.AddAt(SetupStep, &StepCallback{Step: "Setup", Provider: p.Name(), Resource: r.Name(), Callback: r.Setup})
		e.Registry.AddAt(SetupStep, &StepCallback{Step: "Run", Provider: p.Name(), Resource: r.Name(), Callback: r.Run})
		return true
	} else {
		return false
	}
}

func (r *BashResource) Setup(e *RunEnvironment) {
	e.CRDocker.Pull("ubuntu")
}

func (r *BashResource) Run(e *RunEnvironment) {
	e.CRDocker.Run(dockerRunConfig{
		Image:     "ubuntu",
		DestDir:   "/usr/src/myapp",
		SourceDir: Cwd(),
		Cmd:       append([]string{"bash"}, e.Cmd...),
	})
}
