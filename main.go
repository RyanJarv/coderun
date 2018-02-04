package main

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"

	"github.com/rggerst/coderun/coderun"
)

func main() {
	flag.Parse()

	var config = &coderun.ProviderConfig{
		Extension:     path.Ext(flag.Args()[0]),
		Cmd:           flag.Args()[0],
		Args:          flag.Args()[1:len(flag.Args())],
		FullCmdString: strings.Join(flag.Args(), " "),
	}

	provider, err := coderun.CreateProvider(config)
	if err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory")
	}

	var runEnv = &coderun.RunEnvironment{
		Cwd:           cwd,
		Cmd:           flag.Args()[0],
		Args:          flag.Args()[1:len(flag.Args())],
		ArgsString:    strings.Join(flag.Args()[1:len(flag.Args())], " "),
		FullCmdString: strings.Join(flag.Args(), " "),
	}
	provider.Setup(runEnv)
	provider.Run(runEnv)
}
