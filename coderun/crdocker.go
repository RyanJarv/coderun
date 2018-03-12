package coderun

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type CRDocker struct {
	Client  *client.Client
	volumes map[string]string
}

func (d CRDocker) RegisterMount(localPath, dockerPath string) {
	d.volumes[localPath] = dockerPath
}

func (d CRDocker) Pull(image string) {
	Logger.info.Printf("Pulling image: %s", image)
	resp, err := d.Client.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer resp.Close()

	rd := bufio.NewReader(resp)

	rd.WriteTo(os.Stdout)

	Logger.debug.Printf("Done pulling image: %s", image)
}

func (d CRDocker) Run(c dockerRunConfig) {
	ctx := context.Background()

	Logger.info.Printf("Running: %s", c.Cmd)
	config := &container.Config{Image: c.Image, Tty: true, OpenStdin: true, AttachStdin: true, AttachStdout: true, AttachStderr: true}
	if v := c.Cmd; v != nil {
		config.Cmd = v
	}
	if v := c.DestDir; v != "" {
		config.WorkingDir = c.DestDir
	}

	var portBindings nat.PortMap
	port := nat.Port(fmt.Sprintf("%v/tcp", c.Port))

	if c.HostPort == 0 {
		c.HostPort = c.Port
	}

	if c.Port != 0 {
		log.Printf("Setting port %v", c.HostPort)
		portBindings = nat.PortMap{port: []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(c.HostPort)}}}
		config.ExposedPorts = nat.PortSet{
			port: struct{}{},
		}
	}

	hostConfig := &container.HostConfig{PortBindings: portBindings, PidMode: container.PidMode(c.PidMode), Privileged: c.Privileged, NetworkMode: "bridge"}
	if c.NetworkMode != "" {
		hostConfig.NetworkMode = container.NetworkMode(c.NetworkMode)
	}

	if v := c.SourceDir; v != "" {
		Logger.debug.Printf("Sourcedir is %s", v)
		Logger.debug.Printf("Destdir is %s", c.DestDir)
		hostConfig.Mounts = []mount.Mount{{Type: "bind", Source: v, Target: c.DestDir}}
	}
	for l, r := range d.volumes {
		Logger.info.Printf("Attaching bind mount %v to %v", l, r)
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{Type: "bind", Source: l, Target: r})
	}

	netConfig := &network.NetworkingConfig{}

	resp, err := d.Client.ContainerCreate(ctx, config, hostConfig, netConfig, d.newImageName())
	if err != nil {
		panic(err)
	}
	Logger.debug.Printf("Container ID: %s", resp.ID)
	if len(resp.Warnings) > 0 {
		Logger.warn.Printf("Container Warnings: %s", resp.Warnings)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for sig := range ch {
			Logger.error.Printf("Recieved %s, cleaning up", sig.String())
			d.Client.ContainerKill(context.Background(), resp.ID, "SIGTERM")
			timeout := 2 * time.Second
			d.Client.ContainerStop(context.Background(), resp.ID, &timeout)
			d.Client.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{})
		}
	}()

	var errStdout, errStdin error

	hijack, err := d.Client.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{Stream: true, Stdin: true, Stdout: true, Stderr: true})
	if err != nil {
		panic(err)
	}
	defer hijack.Close()

	go func() {
		if c.Stdout != nil {
			_, errStdout = io.Copy(c.Stdout, hijack.Reader)
		} else {
			_, errStdout = io.Copy(os.Stdout, hijack.Reader)
		}
	}()

	if c.Stdin != nil {
		go func() {
			_, errStdin = io.Copy(os.Stdin, hijack.Conn)
		}()
	}

	if errStdout != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}

	if err := d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := d.Client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	d.Client.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{})
}

func (d CRDocker) buildImageStep(source string, args ...string) string {
	var image = d.newImageName()
	var preimage = d.newImageName()
	//append so go will let us pass to a function with a single vervadic parameter
	Exec(append([]string{"/usr/local/bin/docker", "run", "-t", "--name", preimage, "-v", fmt.Sprintf("%s:/usr/local/myapp", Cwd()), "-w", "/usr/local/myapp", source}, args...)...)
	Exec("/usr/local/bin/docker", "commit", preimage, image)
	Exec("/usr/local/bin/docker", "rm", preimage)
	return image
}

func (d CRDocker) Stop(name string) {
	Exec("/usr/local/bin/docker", "stop", name) // Doesn't necessarily stop on it's own
}

func (d CRDocker) getImageName() string {
	image, err := ioutil.ReadFile(".coderun/dockerimage")
	if os.IsNotExist(err) {
		return ""
	} else if err != nil {
		log.Fatal(err)
	}
	return string(image)
}

func (d CRDocker) setImageName(image string) {
	CreateCodeRunDir()
	err := ioutil.WriteFile(".coderun/dockerimage", []byte(image), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (d CRDocker) getOrBuildImage(source string, cmds ...[]string) string {
	var image string
	if image = d.getImageName(); image == "" {
		for _, step := range cmds {
			image = d.buildImageStep(source, step...)
			source = image
		}
		d.setImageName(image)
	}
	return image
}

func (d CRDocker) newImageName() string {
	return fmt.Sprintf("coderun-%s", d.randString())
}

func (d CRDocker) randString() string {
	rand.Seed(rand.NewSource(time.Now().UnixNano()).Int63())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 15)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
