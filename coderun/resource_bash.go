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

func (r *BashResource) Register(e *RunEnvironment) bool {
	return MatchCommandOrExt(e.Cmd, "bash", ".sh")
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
