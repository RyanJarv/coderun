package coderun

import (
	"errors"
	"fmt"
	"github.com/RyanJarv/coderun/coderun/lib"
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
	Resource IResource
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
	Logger.Info.Printf("Adding step %s.%s.%s before %s.%s.%s", s.Provider.Name(), lib.getNameOrEmpty(s.Resource), s.Step, search.Provider, search.Resource, search.Step)
	r.Before = append(r.Before, &StepCallbackSearch{Search: search, Callback: s})
}

func (r *Registry) AddAfter(search *StepSearch, s *StepCallback) {
	Logger.Debug.Printf("Adding step %s.%s.%s after %s.%s.%s", s.Provider.Name(), lib.getNameOrEmpty(s.Resource), s.Step, search.Provider, search.Resource, s.Step)
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
	Logger.Debug.Printf("Searching for callbacks to run before %s.%s.%s", step.Provider.Name(), lib.getNameOrEmpty(step.Resource), step.Step)
	r.runMatching(r.Before, step)
	Logger.Info.Printf("Runlevel %d: Running %s for resource %s, provider %s", r.currentOrder, step.Step, lib.getNameOrEmpty(step.Resource), step.Provider.Name())
	step.Callback(step, step)
	Logger.Debug.Printf("Searching for callbacks to run after %s.%s.%s", step.Provider.Name(), lib.getNameOrEmpty(step.Resource), step.Step)
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
				Logger.Error.Printf("Panic raised while running teardown step. err: %s", err)
			} else if err != nil && r.currentOrder < TeardownStep {
				Logger.Error.Printf("Error occured, cleaning up. err: ", err)
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
		if s.Provider.MatchString(step.Provider.Name()) && s.Resource.MatchString(lib.getNameOrEmpty(step.Resource)) && s.Step.MatchString(step.Step) {
			Logger.Info.Printf("Running %s.%s.%s because it was registered with %s.%s.%s", c.Provider.Name(), lib.getNameOrEmpty(c.Resource), c.Step, step.Provider.Name(), lib.getNameOrEmpty(step.Resource), step.Step)
			c.Callback(c, step)
		}
	}
}
