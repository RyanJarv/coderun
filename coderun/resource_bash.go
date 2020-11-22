package coderun

import (
	"github.com/RyanJarv/coderun/coderun/lib"
	"time"

	"github.com/chzyer/readline"
)

func NewBashResource(r IRunEnvironment) *BashResource {
	return &BashResource{}
}

type BashResource struct {
	IResource
	bash ICRDocker
	env  IRunEnvironment
}

func (r *BashResource) Name() string {
	return "bash"
}

func (r *BashResource) Register(e IRunEnvironment, p IProvider) bool {
	r.env = e
	r.bash = lib.NewCRDocker()
	if s := e.Shell(); s != nil {
		s.AddCompleters(readline.PcItem("bash"))
	}
	if lib.MatchCommandOrExt(e.Cmd(), "bash", ".sh") {
		e.Registry().AddAt(SetupStep, &StepCallback{Step: "Setup", Provider: p, Resource: r, Callback: r.Setup})
		e.Registry().AddAt(SetupStep, &StepCallback{Step: "Run", Provider: p, Resource: r, Callback: r.Run})
		e.Registry().AddAt(TeardownStep, &StepCallback{Step: "Teardown", Provider: p, Resource: r, Callback: r.Teardown})
		return true
	} else {
		return false
	}
}

func (r *BashResource) RegisterMount(local string, remote string) {
	r.bash.RegisterMount(local, remote)
}

func (r *BashResource) Setup(callback *StepCallback, currentStep *StepCallback) {
	r.bash.Pull("bash")
}

func (r *BashResource) Run(callback *StepCallback, currentStep *StepCallback) {
	r.bash.Run(DockerRunConfig{
		Image:     "bash",
		DestDir:   "/usr/src/myapp",
		Attach:    true,
		SourceDir: lib.Cwd(),
		Stdin:     r.env.Stdin(),
		Cmd:       append([]string{"bash"}, r.env.Cmd()...),
	})
}

func (r *BashResource) Teardown(callback *StepCallback, currentStep *StepCallback) {
	r.bash.Teardown(4 * time.Second)
}
