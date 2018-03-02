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

func NewCoderunFs(remotePath string) *CoderunFs {
	log.Printf("Creating new CoderunFs")
	fs := &CoderunFs{
		RemotePath:    remotePath,
		FileResources: map[string]IFileResource{},
	}
	return fs
}

type CoderunFs struct {
	pathfs.FileSystem
	RemotePath    string
	LocalPath     string
	FileResources map[string]IFileResource
}

func (fs *CoderunFs) Setup() {
	os.MkdirAll("/tmp/coderun", 0700)
	tmpdir, err := ioutil.TempDir("/tmp/coderun", "fs")
	if err != nil {
		Logger.error.Fatal(err)
	}
	fs.LocalPath = tmpdir
	nfs := pathfs.NewPathNodeFs(&CoderunFs{FileSystem: pathfs.NewDefaultFileSystem()}, nil)
	server, _, err := nodefs.MountRoot(tmpdir, nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	go server.Serve()
}

func (fs *CoderunFs) AddFileResource(r IFileResource) {
	fs.FileResources[r.Name()] = r
}

func (fs *CoderunFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	if _, ok := fs.FileResources[name]; ok == true {
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

func (fs *CoderunFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	if name == "" {
		dir := make([]fuse.DirEntry, 0, len(fs.FileResources))
		i := 0
		for name, _ := range fs.FileResources {
			i++
			dir[i] = fuse.DirEntry{Name: name, Mode: fuse.S_IFREG}
		}
		return dir, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (fs *CoderunFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	resource, ok := fs.FileResources[name]
	if ok == true {
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
