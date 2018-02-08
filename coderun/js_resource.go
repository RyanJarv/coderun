package coderun

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func JsProvider() *Provider {
	return &Provider{
		RegisterOnCmd: jsRegisterOnCmd,
		Setup:         jsSetup,
		Run:           jsRun,
	}
}

func jsRegisterOnCmd(cmd ...string) bool {
	match, err := regexp.MatchString(`^(([^ ]+/)?node(js)? .*|[\S]+\.js)$`, strings.Join(cmd, " "))
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func jsSetup(r *RunEnvironment) {
	Cmd("/usr/local/bin/docker", "pull", "node")

	if _, err := os.Stat("./package-lock.json"); os.IsNotExist(err) {
		Cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", r.CRDocker.newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", Cwd()), "-w", "/usr/src/myapp", "node", "npm", "install")
	}
}

func jsRun(r *RunEnvironment) {
	r.CRDocker.Run(dockerRunConfig{Image: "node", DestDir: "/usr/src/myapp", SourceDir: Cwd(), Cmd: append([]string{"node"}, r.Cmd...)})
}
