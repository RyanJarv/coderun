package coderun

import "time"

func NewBashResource(r **RunEnvironment) IDockerResource {
	return &BashResource{env: r}
}

type BashResource struct {
	IResource
	bash *CRDocker
	env  **RunEnvironment
}

func (r *BashResource) Name() string {
	return "bash"
}

func (r *BashResource) Register(p IProvider) bool {
	if MatchCommandOrExt((*r.env).Cmd, "bash", ".sh") {
		(*r.env).Registry.AddAt(SetupStep, &StepCallback{Step: "Setup", Provider: p, Resource: r, Callback: r.Setup})
		(*r.env).Registry.AddAt(SetupStep, &StepCallback{Step: "Run", Provider: p, Resource: r, Callback: r.Run})
		(*r.env).Registry.AddAt(TeardownStep, &StepCallback{Step: "Teardown", Provider: p, Resource: r, Callback: r.Teardown})
		return true
	} else {
		return false
	}
}

func (r *BashResource) RegisterMount(local string, remote string) {
	r.bash.RegisterMount(local, remote)
}

func (r *BashResource) Setup(callback *StepCallback, currentStep *StepCallback) {
	r.bash = NewCRDocker()
	r.bash.Pull("bash")
}

func (r *BashResource) Run(callback *StepCallback, currentStep *StepCallback) {
	r.bash.Run(dockerRunConfig{
		Image:     "bash",
		DestDir:   "/usr/src/myapp",
		Attach:    true,
		SourceDir: Cwd(),
		Cmd:       append([]string{"bash"}, (*r.env).Cmd...),
	})
}

func (r *BashResource) Teardown(callback *StepCallback, currentStep *StepCallback) {
	r.bash.Teardown(4 * time.Second)
}
