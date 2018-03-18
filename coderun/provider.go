package coderun

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/RyanJarv/coderun/coderun/shell"
)

type IProvider interface {
	Name() string
	Register(IRunEnvironment) bool
}

type ProviderHookFunc func(IProvider, IRunEnvironment)

type IResource interface {
	Name() string
	Register(IRunEnvironment, IProvider) bool
}

type IRunEnvironment interface {
	Providers() map[string]IProvider
	Registry() *Registry
	Shell() *shell.Shell
	Ask(string) string
	Stdin() io.Reader
	Cmd() []string
}

type RunEnvironment struct {
	Name                string
	EntryPoint          string
	providers           map[string]IProvider
	registeredProviders map[string]IProvider

	cmd         []string
	CodeDir     string
	DependsDir  string
	IgnoreFiles []string
	Flags       map[string]*string
	shell       *shell.Shell
	askR        io.Reader
	askMutex    *sync.Mutex
	stdin       *io.PipeReader
	CRDocker    ICRDocker
	Exec        func(...string) string
	registry    *Registry
}

func (e *RunEnvironment) Providers() map[string]IProvider { return e.providers }
func (e *RunEnvironment) Cmd() []string                   { return e.cmd }
func (e *RunEnvironment) Shell() *shell.Shell             { return e.shell }
func (e *RunEnvironment) Stdin() io.Reader                { return e.stdin }
func (e *RunEnvironment) Registry() *Registry             { return e.registry }
func (e *RunEnvironment) Ask(p string) string {
	e.askMutex.Lock()
	ask := shell.NewAsk(e.askR, os.Stdout)
	fmt.Fprintln(ask, p)
	out, err := bufio.NewReader(ask).ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			Logger.error.Fatal(err)
		}
	}
	e.askMutex.Unlock()
	return string(out)
}

type Stdio struct {
	buf bytes.Buffer
}

type IProviderEnv interface {
}

func CreateRunEnvironment() *RunEnvironment {
	cwd := Cwd()

	ignoreFiles := append(
		readIgnoreFile(path.Join(cwd, ".gitignore")),
		append(
			readIgnoreFile(path.Join(cwd, ".crignore")),
			".coderun",
		)...,
	)

	var runEnv *RunEnvironment
	runEnv = &RunEnvironment{
		providers: map[string]IProvider{
			"mount":  NewMountProvider(runEnv),
			"docker": NewDockerProvider(runEnv),
			"snitch": NewSnitchProvider(runEnv),
			//"lambda": NewAWSLambdaProvider(),
		},
		//Registered: map[string]map[*Provider]*Resource{},
		registeredProviders: map[string]IProvider{},
		CodeDir:             cwd,
		DependsDir:          "",
		IgnoreFiles:         ignoreFiles,
		Name:                path.Base(Cwd()),
		EntryPoint:          "lambda_handler",
		Flags:               make(map[string]*string),
		askMutex:            &sync.Mutex{},
		Exec:                Exec,
		registry:            NewRegistry(),
	}

	var stdinW, askW *io.PipeWriter
	runEnv.stdin, stdinW = io.Pipe()
	runEnv.askR, askW = io.Pipe()
	w := io.MultiWriter(stdinW, askW)
	go func() { _, _ = io.Copy(w, os.Stdin) }()

	return runEnv
}

func run(e *RunEnvironment, cmd []string) {
	e.cmd = cmd
	for _, p := range (*e).providers {
		p.Register(e)
	}
	e.registry.Run()
	Logger.info.Printf("Done running steps")
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
