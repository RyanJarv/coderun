package coderun

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kr/pty"
)

func CreateCodeRunDir() {
	err := os.Mkdir(".coderun", 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func Cmd(c ...string) string {

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

func MatchCommandOrExt(toSearch []string, cmd string, ext string) bool {
	file := filepath.Base(toSearch[0])
	return (file == cmd || filepath.Ext(file) == ext)
}

func Cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory")
	}
	return cwd
}
