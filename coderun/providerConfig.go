package coderun

type ProviderConfig struct {
	Extension     string
	Cmd           string
	Args          []string
	FullCmdString string
}

type RunEnvironment struct {
	Cwd           string
	Cmd           string
	Args          []string
	ArgsString    string
	FullCmdString string
}
