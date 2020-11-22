package lib

type LambdaStatus string

//const (
//	Unknown   ContainerStatus = "Unknown"
//	Created   ContainerStatus = "Created"
//	Running   ContainerStatus = "Running"
//	Removing  ContainerStatus = "Removing"
//	Destroyed ContainerStatus = "Destroyed"
//)
//
//func NewCRDocker() *CRDocker {
//	cli, err := client.NewEnvClient()
//	if err != nil {
//		panic(err)
//	}
//
//	r := &CRDocker{
//		Client:  cli,
//		volumes: map[string]string{},
//		Status:  Unknown,
//	}
//	r.onCtrlC()
//	return r
//}
//
//type CRDocker struct {
//	Client  *client.Client
//	Id      string
//	Info    types.ContainerJSON
//	Status  ContainerStatus
//	hijack  types.HijackedResponse
//	volumes map[string]string
//}
//
//func (d *CRDocker) DockerKillLabel(label string) {
//	args := filters.NewArgs()
//	args.Add("label", label)
//	resp, err := d.Client.ContainerList(context.Background(), types.ContainerListOptions{Filters: args})
//	if err != nil {
//		coderun.Logger.error.Fatal(err)
//	}
//	timeout := time.Second * 4
//	for _, c := range resp {
//		coderun.Logger.info.Printf("Found %s matching label %s", c.ID, label)
//		if err := d.stop(c.ID, timeout); err != nil {
//			coderun.Logger.info.Printf("Could not stop %s in timeout %v, killing", d.Id, timeout)
//			d.kill(c.ID)
//		}
//		d.remove(c.ID)
//	}
//}
//
//func (d *CRDocker) onCtrlC() {
//	ch := make(chan os.Signal, 1)
//	signal.Notify(ch, os.Interrupt)
//	go func() {
//		<-ch
//		d.Teardown(5 * time.Second)
//	}()
//}
//
//func (d *CRDocker) RegisterMount(localPath, dockerPath string) {
//	d.volumes[localPath] = dockerPath
//}
//
//func (d *CRDocker) Pull(image string) {
//	coderun.Logger.info.Printf("Pulling image: %s", image)
//	resp, err := d.Client.ImagePull(context.Background(), image, types.ImagePullOptions{})
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Close()
//
//	rd := bufio.NewReader(resp)
//
//	rd.WriteTo(os.Stdout)
//
//	coderun.Logger.debug.Printf("Done pulling image: %s", image)
//}
//
//func (d *CRDocker) Run(c coderun.dockerRunConfig) {
//	ctx := context.Background()
//
//	coderun.Logger.info.Printf("Running: %s", c.Cmd)
//	config := &container.Config{Image: c.Image, Labels: map[string]string{"coderun": ""}}
//	if c.Attach {
//		config.Tty = true
//		config.OpenStdin = true
//		config.AttachStdin = true
//		config.AttachStdout = true
//		config.AttachStderr = true
//	}
//	if v := c.Cmd; v != nil {
//		config.Cmd = v
//	}
//	if v := c.DestDir; v != "" {
//		config.WorkingDir = c.DestDir
//	}
//
//	var portBindings nat.PortMap
//	port := nat.Port(fmt.Sprintf("%v/tcp", c.Port))
//
//	if c.HostPort == 0 {
//		c.HostPort = c.Port
//	}
//
//	if c.Port != 0 {
//		coderun.Logger.debug.Printf("Setting port %v", c.HostPort)
//		portBindings = nat.PortMap{port: []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: strconv.Itoa(c.HostPort)}}}
//		config.ExposedPorts = nat.PortSet{
//			port: struct{}{},
//		}
//	}
//
//	hostConfig := &container.HostConfig{PortBindings: portBindings, PidMode: container.PidMode(c.PidMode), Privileged: c.Privileged, NetworkMode: "bridge"}
//	if c.NetworkMode != "" {
//		hostConfig.NetworkMode = container.NetworkMode(c.NetworkMode)
//	}
//
//	if v := c.SourceDir; v != "" {
//		coderun.Logger.debug.Printf("Sourcedir is %s", v)
//		coderun.Logger.debug.Printf("Destdir is %s", c.DestDir)
//		hostConfig.Mounts = []mount.Mount{{Type: "bind", Source: v, Target: c.DestDir}}
//	}
//	for l, r := range d.volumes {
//		coderun.Logger.info.Printf("Attaching bind mount %v to %v", l, r)
//		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{Type: "bind", Source: l, Target: r})
//	}
//
//	netConfig := &network.NetworkingConfig{}
//
//	resp, err := d.Client.ContainerCreate(ctx, config, hostConfig, netConfig, d.newImageName())
//	if err != nil {
//		panic(err)
//	}
//	d.Status = Created
//	d.Id = resp.ID
//
//	coderun.Logger.debug.Printf("Container ID: %s", d.Id)
//	if len(resp.Warnings) > 0 {
//		coderun.Logger.warn.Printf("Container Warnings: %s", resp.Warnings)
//	}
//
//	var errStdout, errStdin error
//
//	if c.Attach {
//		coderun.Logger.debug.Printf("Attaching container")
//		d.hijack, err = d.Client.ContainerAttach(ctx, d.Id, types.ContainerAttachOptions{Stream: true, Stdin: true, Stdout: true, Stderr: true})
//		if err != nil {
//			panic(err)
//		}
//
//		go func() {
//			if c.Stdout != nil {
//				_, errStdout = io.Copy(c.Stdout, d.hijack.Reader)
//			} else {
//				_, errStdout = io.Copy(os.Stdout, d.hijack.Reader)
//			}
//		}()
//
//		go func() {
//			if c.Stdin != nil {
//				_, errStdin = io.Copy(d.hijack.Conn, c.Stdin)
//			} else {
//				_, errStdin = io.Copy(d.hijack.Conn, os.Stdin)
//			}
//		}()
//
//		if errStdout != nil {
//			log.Fatal("failed to capture stdout or stderr\n")
//		}
//	}
//
//	if err := d.Client.ContainerStart(ctx, d.Id, types.ContainerStartOptions{}); err != nil {
//		panic(err)
//	}
//	d.Status = Running
//	d.inspect()
//
//	if c.Attach {
//		_, err := d.Client.ContainerWait(ctx, d.Id)
//		if err != nil {
//			coderun.Logger.error.Fatal(err)
//		}
//	}
//}
//
//func (d *CRDocker) inspect() {
//	ctx := context.Background()
//	var err error
//	d.Info, err = d.Client.ContainerInspect(ctx, d.Id)
//	if err != nil {
//		coderun.Logger.error.Fatal(err)
//	}
//}
//
//func (d *CRDocker) Teardown(timeout time.Duration) {
//	if d == nil {
//		coderun.Logger.warn.Printf("CRDocker instance was nil")
//		return
//	}
//	if d.hijack != (types.HijackedResponse{}) {
//		d.hijack.Close()
//	}
//	if d.Status == Destroyed || d.Status == Removing {
//		coderun.Logger.debug.Printf("Container is already in %s state, skipping additional teardown", d.Status)
//		return
//	}
//	d.Status = Removing
//	if err := d.Stop(timeout); err != nil {
//		coderun.Logger.info.Printf("Could not stop %s in timeout %v, killing", d.Id, timeout)
//		d.Kill()
//	}
//	d.Remove()
//	d.Status = Destroyed
//}
//
//func (d *CRDocker) kill(id string) {
//	coderun.Logger.info.Printf("Killing container %s", id)
//	if err := d.Client.ContainerKill(context.Background(), id, "SIGTERM"); err != nil {
//		coderun.Logger.error.Fatal(err)
//	}
//}
//
//func (d *CRDocker) Kill() {
//	d.kill(d.Id)
//}
//
//func (d *CRDocker) stop(id string, timeout time.Duration) error {
//	coderun.Logger.info.Printf("Stopping container %s", id)
//	if err := d.Client.ContainerStop(context.Background(), id, &timeout); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (d *CRDocker) Stop(timeout time.Duration) error {
//	return d.stop(d.Id, timeout)
//}
//
//func (d *CRDocker) remove(id string) {
//	coderun.Logger.info.Printf("Removing container %s", id)
//	if err := d.Client.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{}); err != nil {
//		coderun.Logger.error.Fatal(err)
//	}
//}
//
//func (d *CRDocker) Remove() {
//	d.remove(d.Id)
//}
//
//func (d *CRDocker) buildImageStep(source string, args ...string) string {
//	var image = d.newImageName()
//	var preimage = d.newImageName()
//	//append so go will let us pass to a function with a single vervadic parameter
//	Exec(append([]string{"/usr/local/bin/docker", "run", "-t", "--name", preimage, "-v", fmt.Sprintf("%s:/usr/local/myapp", Cwd()), "-w", "/usr/local/myapp", source}, args...)...)
//	Exec("/usr/local/bin/docker", "commit", preimage, image)
//	Exec("/usr/local/bin/docker", "rm", preimage)
//	return image
//}
//
//func (d *CRDocker) getImageName() string {
//	image, err := ioutil.ReadFile(".coderun/dockerimage")
//	if os.IsNotExist(err) {
//		return ""
//	} else if err != nil {
//		log.Fatal(err)
//	}
//	return string(image)
//}
//
//func (d *CRDocker) setImageName(image string) {
//	CreateCodeRunDir()
//	err := ioutil.WriteFile(".coderun/dockerimage", []byte(image), 0644)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func (d *CRDocker) getOrBuildImage(source string, cmds ...[]string) string {
//	var image string
//	if image = d.getImageName(); image == "" {
//		for _, step := range cmds {
//			image = d.buildImageStep(source, step...)
//			source = image
//		}
//		d.setImageName(image)
//	}
//	return image
//}
//
//func (d *CRDocker) newImageName() string {
//	return fmt.Sprintf("coderun-%s", d.randString())
//}
//
//func (d *CRDocker) randString() string {
//	rand.Seed(rand.NewSource(time.Now().UnixNano()).Int63())
//	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
//	b := make([]rune, 15)
//	for i := range b {
//		b[i] = letterRunes[rand.Intn(len(letterRunes))]
//	}
//	return string(b)
//}
