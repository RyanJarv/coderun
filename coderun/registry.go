package coderun

import "regexp"

const (
	SetupStep = iota * 100
	DeployStep
	RunStep
)

const (
	RunBefore = iota
	RunAfter
)

type Register map[int]StepCallback

type StepCallback struct {
	Callback func(*RunEnvironment, *StepCallback)
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

func NewRegistry() *Registry {
	return &Registry{
		Levels: make([][]*StepCallback, 400),
		Before: []*StepCallbackSearch{},
		After:  []*StepCallbackSearch{},
	}
}

type Registry struct {
	Levels [][]*StepCallback
	Before []*StepCallbackSearch
	After  []*StepCallbackSearch
}

func (r *Registry) AddAt(l int, s *StepCallback) {
	r.Levels[l] = append(r.Levels[l], s)
}

func (r *Registry) AddBefore(search *StepSearch, s *StepCallback) {
	r.Before = append(r.Before, &StepCallbackSearch{Search: search, Callback: s})
}

func (r *Registry) AddAfter(search *StepSearch, s *StepCallback) {
	r.After = append(r.After, &StepCallbackSearch{Search: search, Callback: s})
}

func (r *Registry) Run(e *RunEnvironment) {
	for order, steps := range r.Levels {
		for _, step := range steps {
			r.runMatching(e, r.Before, step)
			Logger.info.Printf("Runlevel %d: Running %s for resource %s, provider %s", order, step.Step, step.Resource, step.Provider)
			step.Callback(e, step)
			r.runMatching(e, r.After, step)
		}
	}
}

func (r *Registry) runMatching(e *RunEnvironment, callbackSearchList []*StepCallbackSearch, step *StepCallback) {
	for _, callbackSearch := range callbackSearchList {
		c := callbackSearch.Callback
		s := callbackSearch.Search
		if s.Provider.MatchString(step.Provider.Name()) && s.Resource.MatchString(step.Resource.Name()) && s.Step.MatchString(step.Step) {
			Logger.info.Printf("Running %s.%s.%s because it was registered with %s.%s.%s", c.Provider, c.Resource, c.Step, step.Provider, step.Resource, step.Step)
			c.Callback(e, c)
		}
	}
}
