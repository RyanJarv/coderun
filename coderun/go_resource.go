package coderun

import (
	"log"
	"os"
	"regexp"
	"strings"
)

func GoProvider() *Provider {
	return &Provider{
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

func goSetup(r *RunEnvironment) {
	r.CRDocker.Pull("golang")

	if _, err := os.Stat("./glide.lock"); os.IsNotExist(err) {
		r.CRDocker.Run(dockerRunConfig{Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: Cwd(), Cmd: []string{"sh", "-c", "curl https://glide.sh/get | sh && glide init --non-interactive && glide install"}})
	}
}

func goRun(r *RunEnvironment) {
	r.CRDocker.Run(dockerRunConfig{Image: "golang", DestDir: "/go/src/localhost/myapp", SourceDir: Cwd(), Cmd: []string{"go", "run"}})
}
