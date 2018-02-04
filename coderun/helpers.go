package coderun

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"

	"github.com/kr/pty"
)

func getOrBuildImage(source string, cmds ...[]string) string {
	var image string
	if image = getImageName(); image == "" {
		for _, step := range cmds {
			image = buildImageStep(source, step...)
			source = image
		}
		setImageName(image)
	}
	return image
}

func buildImageStep(source string, args ...string) string {
	var image = newImageName()
	var preimage = newImageName()
	//append so go will let us pass to a function with a single vervadic parameter
	cmd(append([]string{"/usr/local/bin/docker", "run", "-t", "--name", preimage, "-v", fmt.Sprintf("%s:/usr/local/myapp", cwd()), "-w", "/usr/local/myapp", source}, args...)...)
	cmd("/usr/local/bin/docker", "commit", preimage, image)
	cmd("/usr/local/bin/docker", "rm", preimage)
	return image
}

func dockerPull(image string) {
	cmd("/usr/local/bin/docker", "pull", image)
}

func dockerRun(image string, port int, path string, args ...string) {
	name := newImageName()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Recieved %s, cleaning up", sig.String())
			dockerStop(name)
		}
	}()
	p := strconv.Itoa(port)
	cmd(append([]string{"/usr/local/bin/docker", "run", "-e", p, "-p", fmt.Sprintf("%s:%s", p, p), "-it", "--rm", "--name", name, "-v", fmt.Sprintf("%s:%s", cwd(), path), "-w", path, image}, args...)...)
}

func dockerStop(name string) {
	cmd("/usr/local/bin/docker", "stop", name) // Doesn't necessarily stop on it's own
}

func getImageName() string {
	image, err := ioutil.ReadFile(".coderun/dockerimage")
	if os.IsNotExist(err) {
		return ""
	} else if err != nil {
		log.Fatal(err)
	}
	return string(image)
}

func setImageName(image string) {
	createCodeRunDir()
	err := ioutil.WriteFile(".coderun/dockerimage", []byte(image), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func createCodeRunDir() {
	err := os.Mkdir(".coderun", 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func cmd(c ...string) string {

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.Command(c[0], c[1:]...)
	log.Printf("%v", cmd.Args)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	log.Printf("Running command and waiting for it to finish...")
	tty, err := pty.Start(cmd)
	if err != nil {
		fmt.Println("Error start cmd", err)
	}
	defer tty.Close()

	go func() {
		scanner := bufio.NewScanner(tty)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	go func() {
		io.Copy(tty, os.Stdin)
	}()

	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	log.Printf("Done with command %s", cmd.Args)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}

	outStr := string(stdoutBuf.Bytes())

	return outStr
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
	return string(b)
}

func cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory")
	}
	return cwd
}
