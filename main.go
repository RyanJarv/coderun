package main

import (
	"flag"
	L "github.com/RyanJarv/coderun/coderun/logger"
	"log"

	"github.com/RyanJarv/coderun/coderun"
)

func main() {
	runEnv := coderun.CreateRunEnvironment()
	flag.Parse()

	env, err := coderun.Setup(runEnv, flag.Args())
	L.Debug.Printf("Run Environment: %v:", env)
	if err != nil {
		log.Fatal(err)
	}
}
