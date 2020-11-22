package coderun

import (
	"bytes"
	"github.com/RyanJarv/coderun/coderun/lib"
	"os"
	"path"
	"strings"

	"github.com/RyanJarv/coderun/coderun/shell"
)

type IProvider interface {
	Name() string
	Register(IRunEnvironment) bool
	Setup(*StepCallback, *StepCallback) // callback, currentStep
	Run(*StepCallback, *StepCallback) // callback, currentStep
}

type Provider struct {
	name string
	resources           map[string]ILambdaResource
	registeredResources map[string]ILambdaResource
}

func (p *Provider) Name() string 							{ return p.name }
func (p *Provider) Setup(_ *StepCallback, _ *StepCallback)  { }
func (p *Provider) Run(_ *StepCallback, _ *StepCallback) 	{ }

type ProviderHookFunc func(IProvider, IRunEnvironment)

type IResource interface {
	Name() string
	Register(IRunEnvironment, IProvider) bool
	Setup(*StepCallback, *StepCallback) // callback, currentStep
	Run(*StepCallback, *StepCallback) // callback, currentStep
}

type Resource struct {
	name string
}
func (p *Resource) Name() string 							{ return p.name }
func (p *Resource) Setup(_ *StepCallback, _ *StepCallback)  { }
func (p *Resource) Run(_ *StepCallback, _ *StepCallback) 	{ }

type IRunEnvironment interface {
	Providers() map[string]IProvider
	Registry() IRegistry
	Shell() *shell.Shell
	Stdin() *shell.StdinSwitch
	Cmd() []string
	DependsDir() string
	CodeDir() string
	SetCmd([]string)
	Docker() ICRDocker
	IgnoreFiles() []string
}

type RunEnvironment struct {
	Name                string
	EntryPoint          string
	providers           map[string]IProvider
	registeredProviders map[string]IProvider

	cmd         []string
	codeDir     string
	dependsDir  string
	ignoreFiles []string
	Flags       map[string]*string
	shell       *shell.Shell
	stdin       *shell.StdinSwitch
	CRDocker    ICRDocker
	Exec        func(...string) string
	registry    IRegistry
}

func (e *RunEnvironment) Providers() map[string]IProvider { return e.providers }
func (e *RunEnvironment) SetCmd(c []string)               { e.cmd = c }
func (e *RunEnvironment) Cmd() []string                   { return e.cmd }
func (e *RunEnvironment) CodeDir() string                 { return e.codeDir }
func (e *RunEnvironment) DependsDir() string              { return e.dependsDir }
func (e *RunEnvironment) Shell() *shell.Shell             { return e.shell }
func (e *RunEnvironment) Stdin() *shell.StdinSwitch       { return e.stdin }
func (e *RunEnvironment) Registry() IRegistry             { return e.registry }
func (e *RunEnvironment) Docker() ICRDocker               { return e.CRDocker }
func (e *RunEnvironment) IgnoreFiles() []string           { return e.ignoreFiles }

type Stdio struct {
	buf bytes.Buffer
}

//type IProviderEnv interface {
//}

func CreateRunEnvironment() *RunEnvironment {
	cwd := lib.Cwd()

	ignoreFiles := append(
		lib.readIgnoreFile(path.Join(cwd, ".gitignore")),
		append(
			lib.readIgnoreFile(path.Join(cwd, ".crignore")),
			".coderun",
		)...,
	)

	var runEnv *RunEnvironment
	runEnv = &RunEnvironment{
		providers: map[string]IProvider{
			//"mount":  NewMountProvider(runEnv),
			//"docker": NewDockerProvider(runEnv),
			//"snitch": NewSnitchProvider(runEnv),
			"lambda": NewLambdaProvider(runEnv),
		},
		//Registered: map[string]map[*Provider]*Resource{},
		registeredProviders: map[string]IProvider{},
		codeDir:             cwd,
		ignoreFiles:         ignoreFiles,
		Name:                path.Base(lib.Cwd()),
		EntryPoint:          "lambda_handler",
		Flags:               make(map[string]*string),
		stdin:               shell.NewStdinSwitch(os.Stdin, os.Stdout),
		Exec:                lib.Exec,
	}

	return runEnv
}

func run(e *RunEnvironment, cmd []string) {
	e.registry = NewRegistry()
	e.SetCmd(cmd)
	for _, p := range e.providers {
		p.Register(e)
	}
	e.registry.Run()
	Logger.Info.Printf("Done running steps")
}

func Setup(e *RunEnvironment, cmd []string) (IRunEnvironment, error) {
	if len(cmd) == 0 {
		e.shell = shell.NewShell()
		e.shell.Start(func(line string) {
			cmd = strings.Split(line, " ")
			run(e, cmd)
		})
	} else {
		run(e, cmd)
	}

	return e, nil
}
func Deploy(e *RunEnvironment, cmd []string) (IRunEnvironment, error) {
	if len(cmd) == 0 {
		e.shell = shell.NewShell()
		e.shell.Start(func(line string) {
			cmd = strings.Split(line, " ")
			run(e, cmd)
		})
	} else {
		run(e, cmd)
	}

	return e, nil
}
func Run(e *RunEnvironment, cmd []string) (IRunEnvironment, error) {
	if len(cmd) == 0 {
		e.shell = shell.NewShell()
		e.shell.Start(func(line string) {
			cmd = strings.Split(line, " ")
			run(e, cmd)
		})
	} else {
		run(e, cmd)
	}

	return e, nil
}
