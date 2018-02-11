package main

import (
	"flag"
	"log"

	"github.com/ryanjarv/coderun/coderun"
)

func main() {
	runEnv := coderun.CreateRunEnvironment()
	runEnv.Flags["provider"] = flag.String("p", "", "Use given provider (docker|lambda)")
	flag.Parse()
	runEnv.Cmd = flag.Args()

	_, err := coderun.Setup(runEnv)
	if err != nil {
		log.Fatal(err)
	}
}
