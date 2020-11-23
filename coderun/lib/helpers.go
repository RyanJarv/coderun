package lib

import (
	"bufio"
	"bytes"
	"fmt"
	L "github.com/RyanJarv/coderun/coderun/logger"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kr/pty"
)

func CreateCodeRunDir() {
	err := os.Mkdir(".coderun", 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func Exec(c ...string) string {

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.Command(c[0], c[1:]...)
	L.Info.Printf("Command arguments: %v", cmd.Args)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	L.Info.Printf("Running command and waiting for it to finish...")
	tty, err := pty.Start(cmd)
	if err != nil {
		L.Error.Fatal("Error start cmd", err)
	}
	defer tty.Close()

	go func() {
		scanner := bufio.NewScanner(tty)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	if err != nil {
		L.Error.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	L.Info.Printf("Done with command %s", cmd.Args)
	if err != nil {
		L.Error.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		L.Error.Fatal("failed to capture stdout or stderr\n")
	}

	outStr := string(stdoutBuf.Bytes())

	return outStr
}

func MatchCommandOrExt(toSearch []string, cmd string, ext string) bool {
	L.Debug.Printf("MatchCommandOrExt got command: %s", strings.Join(toSearch, " "))
	if len(toSearch) == 0 {
		return false
	}
	file := filepath.Base(toSearch[0])
	return (file == cmd || filepath.Ext(file) == ext)
}

func Cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		L.Error.Fatal("Error getting current working directory")
	}
	return cwd
}

func RandString(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ReadIgnoreFile(f string) []string {
	var ignoreFiles []string
	file, err := ioutil.ReadFile(f)
	if os.IsNotExist(err) {
		ignoreFiles = []string{}
	} else if err != nil {
		log.Fatal(err)
	} else {
		ignoreFiles := make([]string, len(file))
		for i, l := range file {
			ignoreFiles[i] = strings.Trim(string(l), " \t")
		}
	}
	return ignoreFiles
}

type INameable interface {
	Name() string
}

func GetNameOrEmpty(s INameable) string {
	n := ""
	if s != nil {
		n = s.Name() //Resource can be nil
	}
	return n
}
