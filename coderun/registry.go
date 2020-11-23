package coderun

import (
	"errors"
	"fmt"
	"github.com/RyanJarv/coderun/coderun/lib"
	L "github.com/RyanJarv/coderun/coderun/logger"
	"os"
	"os/signal"
	"regexp"
)

const (
	SetupStep = iota*100 + 100
	DeployStep
	RunStep
	TeardownStep
)

const (
	RunBefore = iota
	RunAfter
)

type Register map[int]StepCallback

type StepCallback struct {
	Callback func(*StepCallback, *StepCallback)
	Provider IProvider
	Step     string
}

type StepCallbackSearch struct {
	Callback *StepCallback
	Search   *StepSearch
}

type StepSearch struct {
	Provider *regexp.Regexp
	Resource *regexp.Regexp
	Step     *regexp.Regexp
}

func NewRegistry() IRegistry {
	r := &Registry{
		Levels:   make([][]*StepCallback, 499),
		Before:   []*StepCallbackSearch{},
		After:    []*StepCallbackSearch{},
		Teardown: false,
	}
	r.onCtrlC()
	return r
}

type IRegistry interface {
	AddAt(int, *StepCallback)
	AddBefore(*StepSearch, *StepCallback)
	AddAfter(*StepSearch, *StepCallback)
	Run()
}

type Registry struct {
	Levels       [][]*StepCallback
	currentOrder int
	Before       []*StepCallbackSearch
	After        []*StepCallbackSearch
	Teardown     bool
}

func (r *Registry) AddAt(l int, s *StepCallback) {
	r.Levels[l] = append(r.Levels[l], s)
}

func (r *Registry) AddBefore(search *StepSearch, s *StepCallback) {
	L.Info.Printf("Adding step %s.%s.%s before %s.%s.%s", s.Provider.Name(), lib.GetNameOrEmpty(s.Provider), s.Step, search.Provider, search.Resource, search.Step)
	r.Before = append(r.Before, &StepCallbackSearch{Search: search, Callback: s})
}

func (r *Registry) AddAfter(search *StepSearch, s *StepCallback) {
	L.Debug.Printf("Adding step %s.%s.%s after %s.%s.%s", s.Provider.Name(), lib.GetNameOrEmpty(s.Provider), s.Step, search.Provider, search.Resource, s.Step)
	r.After = append(r.After, &StepCallbackSearch{Search: search, Callback: s})
}

func (r *Registry) onCtrlC() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		// Only run teardown steps after this is triggered
		r.Teardown = true
	}()
}

func (r *Registry) runStep(step *StepCallback) (err error) {
	defer func() {
		if recov := recover(); recov != nil {
			err = errors.New(fmt.Sprint(recov))
		}
	}()
	L.Debug.Printf("Searching for callbacks to run before %s.%s", step.Provider.Name(), step.Step)
	r.runMatching(r.Before, step)
	L.Info.Printf("Runlevel %d: Running %s for provider %s", r.currentOrder, step.Step, step.Provider.Name())
	step.Callback(step, step)
	L.Debug.Printf("Searching for callbacks to run after %s.%s", step.Provider.Name(), step.Step)
	r.runMatching(r.After, step)
	return
}

func (r *Registry) run() {
	var steps []*StepCallback
	for r.currentOrder, steps = range r.Levels {
		if r.Teardown && r.currentOrder < TeardownStep {
			continue
		}
		for _, step := range steps {
			err := r.runStep(step)
			if err != nil && r.currentOrder >= TeardownStep {
				L.Error.Printf("Panic raised while running teardown step. err: %s", err)
			} else if err != nil && r.currentOrder < TeardownStep {
				L.Error.Printf("Error occured, cleaning up. err: ", err)
				r.Teardown = true
			}
		}
	}
}

func (r *Registry) Run() {
	defer func() {
	}()
	r.run()
}

func (r *Registry) runMatching(callbackSearchList []*StepCallbackSearch, step *StepCallback) {
	for _, callbackSearch := range callbackSearchList {
		c := callbackSearch.Callback
		s := callbackSearch.Search
		if s.Provider.MatchString(step.Provider.Name()) && s.Step.MatchString(step.Step) {
			L.Info.Printf("Running %s.%s.%s because it was registered with %s.%s", c.Provider.Name(), c.Step, step.Provider.Name(), step.Step)
			c.Callback(c, step)
		}
	}
}
