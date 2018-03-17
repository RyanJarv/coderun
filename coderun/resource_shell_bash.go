package coderun

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/RyanJarv/coderun/coderun/prompt"
	"github.com/chzyer/readline"
	"golang.org/x/crypto/ssh/terminal"
)

func NewShellBashResource(e IRunEnvironment) IShellResource {
	return &ShellBashResource{
		registry: NewRegistry(),
	}
}

type ShellBashResource struct {
	IResource
	bash                *CRDocker
	shell               *exec.Cmd
	tty                 *os.File
	stdin               *io.PipeWriter
	stdout              *io.PipeReader
	registry            *Registry
	oldstate            *terminal.State
	input               chan byte
	providers           map[string]IProvider
	RegisteredProviders map[string]IProvider
}

func (r *ShellBashResource) Name() string {
	return "shellbash"
}

func (r *ShellBashResource) Register(e IRunEnvironment, p IProvider) bool {
	r.providers = e.Providers()
	cmd := e.Cmd()
	if len(cmd) == 0 || len(cmd) == 1 && cmd[0] == "bash" {
		e.Registry().AddAt(SetupStep, &StepCallback{Step: "Setup", Provider: p, Resource: r, Callback: r.Setup})
		//e.Registry().AddAt(RunStep, &StepCallback{Step: "Run", Provider: p, Resource: r, Callback: r.Run})
		e.Registry().AddAt(TeardownStep, &StepCallback{Step: "Teardown", Provider: p, Resource: r, Callback: r.Teardown})
		return true
	} else {
		return false
	}
}

func (r *ShellBashResource) Setup(callback *StepCallback, currentStep *StepCallback) {
	p := prompt.NewPrompt()
	p.AddCompleters(
		readline.PcItem("test"),
	)
	p.Start()
}

func (r *ShellBashResource) Run(callback *StepCallback, currentStep *StepCallback) {
	//runShellCmds(r.shell, r.tty, r.NewCmd)
}

func (r *ShellBashResource) Teardown(callback *StepCallback, currentStep *StepCallback) {
	Logger.info.Printf("teardown")
	r.tty.Close()
	r.shell.Process.Kill()
}

func (r *ShellBashResource) NewCmd(cmd []string) []string {
	// Use copied cmdEnv and check for providers on each cmd
	env := &RunEnvironment{
		cmd:       cmd,
		providers: r.providers,
		registry:  r.registry,
	}
	Logger.debug.Printf("Got command: %s", strings.Join(cmd, " "))

	registered := false
	for n, bp := range r.providers {
		Logger.debug.Printf("Calling register on Provider %s", bp.Name())
		if bp.Register(env) {
			Logger.debug.Printf("Registered %s", bp.Name())
			registered = true
			if n == "docker" {
				cmd = []string{}
			}
		}
	}

	if registered {
		r.registry.Run()
		Logger.info.Printf("Done running cmd steps")
	}
	return cmd
}
