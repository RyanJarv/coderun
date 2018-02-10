package coderun

type RegisterFunc func(RunEnvironment, IProviderEnv) bool
type RunFunc func(RunEnvironment, IProviderEnv)
type SetupFunc func(RunEnvironment, IProviderEnv)

type Resource struct {
	Name        string
	Register    RegisterFunc
	Setup       SetupFunc
	Run         RunFunc
	ProviderEnv IProviderEnv
}
