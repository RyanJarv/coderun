package main

import (
	"flag"
	"log"

	"github.com/ryanjarv/coderun/coderun"
)

func main() {
	flag.Parse()

	var config = &coderun.ProviderConfig{
		Cmd: flag.Args()[0],
	}

	_, err := coderun.GetProvider(config)
	if err != nil {
		log.Fatal(err)
	}
}
