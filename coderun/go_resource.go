package coderun

import (
	"log"
	"os"
	"regexp"
	"strings"
)

func GoResource() *Resource {
	return &Resource{
		RegisterOnCmd: goRegisterOnCmd,
		Setup:         goSetup,
		Run:           goRun,
	}
}

func goRegisterOnCmd(cmd ...string) bool {
	match, err := regexp.MatchString(`^(([^ ]+/)?go .*|[\S]+\.go)$`, strings.Join(cmd, " "))
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func goSetup(r IRunEnvironment) {
	r.(RunEnvironment).CRDocker.(ICRDocker).Pull("golang")

	if _, err := os.Stat("./glide.lock"); os.IsNotExist(err) {
		r.(RunEnvironment).CRDocker.(CRDocker).Run(dockerRunConfig{Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: Cwd(), Cmd: []string{"sh", "-c", "curl https://glide.sh/get | sh && glide init --non-interactive && glide install"}})
	}
}

func goRun(r IRunEnvironment) {
	r.(RunEnvironment).CRDocker.(CRDocker).Run(dockerRunConfig{Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: Cwd(), Cmd: []string{"go", "run"}})
}
