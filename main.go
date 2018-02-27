package main

import (
	"flag"
	"log"

	"github.com/RyanJarv/coderun/coderun"
)

func main() {
	runEnv := coderun.CreateRunEnvironment()
	runEnv.Flags["provider"] = flag.String("p", "", "Use given provider (docker|lambda)")
	logLevel := flag.String("l", "error", "Set log level (debug|info|warn|error)")
	flag.Parse()
	runEnv.Cmd = flag.Args()

	coderun.Logger = coderun.SetupLogger(*logLevel)

	_, err := coderun.Setup(runEnv)
	if err != nil {
		log.Fatal(err)
	}
}
