package coderun

type RegisterFunc func(*RunEnvironment, IProviderEnv) bool
type RunFunc func(*RunEnvironment, IProviderEnv)
type SetupFunc func(*RunEnvironment, IProviderEnv)
type DeployFunc func(*RunEnvironment, IProviderEnv)

type Resource struct {
	Name        string
	Register    RegisterFunc
	Setup       SetupFunc
	Deploy      DeployFunc
	Run         RunFunc
	ProviderEnv IProviderEnv
}
