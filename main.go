package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func main() {
	flag.Parse()
	var file = flag.Arg(0)
	var ext = path.Ext(flag.Arg(0))

	d, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory")
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	switch ext {
	case ".py":
		setupPython()
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", d), "-w", "/usr/src/myapp", "python:3", "sh", "-c", fmt.Sprintf(". .coderun/venv/bin/activate && python %s", file))
	case ".rb":
		setupRuby()
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", d), "-w", "/usr/src/myapp", "ruby:2.1", "ruby", file)
	case ".go":
		setupGo()
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/go/src/localhost/myapp", d), "-w", "/go/src/localhost/myapp", "golang", "go", "run", file)
	default:
		cmd("echo", "default")
	}
}

func setupPython() {
	d, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory")
	}

	cmd("/usr/local/bin/docker", "pull", "python:3")

	if _, err := os.Stat("./requirements.txt"); err == nil {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", d), "-w", "/usr/src/myapp", "python:3", "python", "-m", "venv", ".coderun/venv")
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/usr/src/myapp", d), "-w", "/usr/src/myapp", "python:3", "sh", "-c", ". .coderun/venv/bin/activate && pip install -r ./requirements.txt")
	}
}

func setupRuby() {
	d, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory")
	}

	cmd("/usr/local/bin/docker", "pull", "ruby:2.1")

	if _, err := os.Stat("./Gemfile"); err == nil {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", "my-running-script1", "-v", fmt.Sprintf("%s:/usr/src/myapp", d), "-w", "/usr/src/myapp", "ruby:2.1", "bundler", "install", "--path", ".coderun/vendor/bundle")
	}
}

func setupGo() {
	d, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory")
	}

	cmd("/usr/local/bin/docker", "pull", "golang")

	var image string
	if image, _ = getImageName(); image == "" {
		var image = newImageName()
		var preimage = newImageName()
		cmd("/usr/local/bin/docker", "run", "-t", "--name", preimage, "-v", fmt.Sprintf("%s:/go/src/localhost/myapp", d), "-w", "/go/src/localhost/myapp", "golang", "sh", "-c", "curl https://glide.sh/get | sh")
		cmd("/usr/local/bin/docker", "commit", preimage, image)
		cmd("/usr/local/bin/docker", "rm", preimage)
		err = setImageName(image)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("./glide.lock"); os.IsNotExist(err) {
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/go/src/localhost/myapp", d), "-w", "/go/src/localhost/myapp", image, "glide", "init", "--non-interactive")
		cmd("/usr/local/bin/docker", "run", "-t", "--rm", "--name", newImageName(), "-v", fmt.Sprintf("%s:/go/src/localhost/myapp", d), "-w", "/go/src/localhost/myapp", image, "glide", "install")
	}
}

func getImageName() (string, error) {
	var image, err = ioutil.ReadFile(".coderun/dockerimage")
	if err != nil {
		return "", err
	}
	return string(image), nil
}

func setImageName(image string) error {
	return ioutil.WriteFile(".coderun/dockerimage", []byte(image), 0644)
}

func cmd(c ...string) string {
	cmd := exec.Command(c[0], c[1:]...)
	log.Printf("%v", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("output: %s", out.String())
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: %s", err))
	}
	return out.String()
}

func newImageName() string {
	return fmt.Sprintf("coderun-%s", randString())
}

func randString() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 15)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	log.Printf("randString: %s", string(b))
	return string(b)
}
