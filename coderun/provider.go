package coderun

import (
	"github.com/RyanJarv/coderun/coderun/lib"
	L "github.com/RyanJarv/coderun/coderun/logger"
	"github.com/RyanJarv/coderun/coderun/shell"
	"os"
	"path"
	"strings"
)

type IRunEnvironment interface {
	Registry() IRegistry
	Shell() *shell.Shell
	Stdin() *shell.StdinSwitch
	Cmd() []string
}

type IProvider interface {
	Name() string
	Register(IRunEnvironment, IProvider) bool
	Run(*StepCallback, *StepCallback) // callback, currentStep
}

type RunEnvironment struct {
	Name                string
	EntryPoint          string
	providers           map[string]IProvider
	registeredProviders map[string]IProvider

	cmd         []string
	codeDir     string
	ignoreFiles []string
	shell       *shell.Shell
	stdin       *shell.StdinSwitch
	CRDocker    ICRDocker
	Exec        func(...string) string
	registry    IRegistry
}


func (e *RunEnvironment) SetCmd(c []string)               { e.cmd = c }
func (e *RunEnvironment) Cmd() []string                   { return e.cmd }
func (e *RunEnvironment) Shell() *shell.Shell             { return e.shell }
func (e *RunEnvironment) Stdin() *shell.StdinSwitch       { return e.stdin }
func (e *RunEnvironment) Registry() IRegistry             { return e.registry }

func CreateRunEnvironment() *RunEnvironment {
	cwd := lib.Cwd()

	ignoreFiles := append(
		lib.ReadIgnoreFile(path.Join(cwd, ".gitignore")),
		append(
			lib.ReadIgnoreFile(path.Join(cwd, ".crignore")),
			".coderun",
		)...,
	)

	var runEnv *RunEnvironment
	return &RunEnvironment{
		providers: map[string]IProvider{
			"docker": NewDockerProvider(runEnv),
		},
		registeredProviders: map[string]IProvider{},
		ignoreFiles:         ignoreFiles,
		Name:                path.Base(lib.Cwd()),
		stdin:               shell.NewStdinSwitch(os.Stdin, os.Stdout),
		Exec:                lib.Exec,
	}
}

func run(e *RunEnvironment, cmd []string) {
	e.registry = NewRegistry()
	e.SetCmd(cmd)
	for _, p := range e.providers {
		p.Register(e, nil)
	}
	e.registry.Run()
	L.Info.Printf("Done running steps")
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
