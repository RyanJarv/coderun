package main

import (
	"flag"
	"log"
	"os"
	"path"

	"./coderun"
)

func main() {
	flag.Parse()
	var file = flag.Arg(0)
	var ext = path.Ext(flag.Arg(0))

	var config = &coderun.ProviderConfig{
		Extension: ext,
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
		FilePath: file,
		Cwd:      cwd,
	}
	provider.Setup(runEnv)
	provider.Run(runEnv)
}
