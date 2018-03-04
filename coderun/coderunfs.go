package coderun

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

func NewCoderunFs(remotePath string) CoderunFs {
	os.MkdirAll("/tmp/coderun", 0700)

	tmpdir, err := ioutil.TempDir("/tmp/coderun", "fs")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("CoderunFs local path is %s", tmpdir)

	crfs := CoderunFs{
		FileSystem:    pathfs.NewDefaultFileSystem(),
		remotePath:    remotePath,
		localPath:     tmpdir,
		fileResources: map[string]IFileResource{},
	}

	crfs.nfs = pathfs.NewPathNodeFs(crfs, nil)

	crfs.server, _, err = nodefs.MountRoot(tmpdir, crfs.nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}

	return crfs
}

type CoderunFs struct {
	pathfs.FileSystem
	remotePath    string
	localPath     string
	fileResources map[string]IFileResource
	server        *fuse.Server
	nfs           *pathfs.PathNodeFs
}

func (fs CoderunFs) Setup() {
}

func (fs CoderunFs) Serve() {
	Logger.debug.Printf("Running CoderunFs.Serve")
	fs.server.Serve()
}

func (fs CoderunFs) AddFileResource(r IFileResource) {
	Logger.debug.Printf("Running CoderunFs.AddFileResource with %s", r.Path())
	fs.fileResources[r.Path()] = r
	for n, v := range fs.fileResources {
		log.Printf("AddFileResource File Exists: %s = %v", n, v)
	}
}

func (fs CoderunFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	Logger.debug.Printf("Running CoderunFs.GetAdttr with %s", name)
	if _, ok := fs.fileResources[name]; ok == true {
		Logger.debug.Printf("Found file %s", name)
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
		}, fuse.OK

	} else if name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (fs CoderunFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	Logger.debug.Printf("Running CoderunFs.opendir with %s", name)
	if name == "" {
		dir := make([]fuse.DirEntry, 0, len(fs.fileResources))
		i := 0
		for name, _ := range fs.fileResources {
			i++
			dir[i] = fuse.DirEntry{Name: name, Mode: fuse.S_IFREG}
		}
		return dir, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (fs CoderunFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	Logger.debug.Printf("Running CoderunFs.open with %s", name)
	resource, ok := fs.fileResources[name]
	if ok == false {
		return nil, fuse.ENOENT
	}
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}
	handle := resource.Open()
	buf := new(bytes.Buffer)
	buf.ReadFrom(handle)
	return nodefs.NewDataFile(buf.Bytes()), fuse.OK
}

//fs needs to be a pointer because it gets registered as a step before setup is done
func (fs *CoderunFs) ConnectDocker(s *StepCallback, step *StepCallback) {
	Logger.info.Printf("fs: %v", fs)
	Logger.info.Printf("localPath: %v, remotePath: %v", fs.localPath, fs.remotePath)
	step.Resource.(IDockerResource).RegisterMount(fs.localPath, fs.remotePath)
}
