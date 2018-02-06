package coderun

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
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

func dockerPull(cli *client.Client, image string) {
	log.Printf("Pulling image: %s", image)
	_, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("Done pulling image: %s", image)
}

type dockerRunConfig struct {
	Client    *client.Client
	Image     string
	Port      int
	Cmd       []string
	SourceDir string
	DestDir   string
	Mounts    mount.Mount
}

func dockerRun(c dockerRunConfig) {
	ctx := context.Background()
	cli := c.Client
	port := nat.Port(fmt.Sprintf("%v/tcp", c.Port))

	var portBindings nat.PortMap
	var exposedPorts nat.PortSet
	if c.Port != 0 {
		portBindings = nat.PortMap{port: []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: strconv.Itoa(c.Port)}}}
		exposedPorts = nat.PortSet{
			port: struct{}{},
		}
	}

	m := []mount.Mount{{Type: "bind", Source: c.SourceDir, Target: c.DestDir}}

	log.Printf("Commands: %s", c.Cmd)
	log.Printf("Bindings: %v", portBindings)
	resp, err := cli.ContainerCreate(ctx, &container.Config{Image: c.Image, Cmd: c.Cmd, WorkingDir: c.DestDir, ExposedPorts: exposedPorts, Tty: true, OpenStdin: true, AttachStdin: true, AttachStdout: true, AttachStderr: true}, &container.HostConfig{Mounts: m, PortBindings: portBindings}, &network.NetworkingConfig{}, newImageName())
	if err != nil {
		panic(err)
	}
	log.Printf("Container ID: %s", resp.ID)
	log.Printf("Container Warnings: %s", resp.Warnings)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for sig := range ch {
			log.Printf("Recieved %s, cleaning up", sig.String())
			c.Client.ContainerKill(context.Background(), resp.ID, "SIGTERM")
			timeout := 2 * time.Second
			c.Client.ContainerStop(context.Background(), resp.ID, &timeout)
			cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{})
		}
	}()

	var errStdout error

	hijack, err := cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{Stream: true, Stdin: true, Stdout: true, Stderr: true})
	if err != nil {
		panic(err)
	}
	defer hijack.Close()

	go func() {
		_, errStdout = io.Copy(os.Stdout, hijack.Reader)
	}()

	if errStdout != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{})
}

func dockerStop(name string) {
	cmd("/usr/local/bin/docker", "stop", name) // Doesn't necessarily stop on it's own
}

type dockerStopConfig struct {
	Client  *client.Client
	ID      string
	Timeout *time.Duration
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
