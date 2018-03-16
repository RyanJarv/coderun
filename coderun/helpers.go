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
	"path/filepath"
	"strings"
	"syscall"

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
	Logger.info.Printf("Command arguments: %v", cmd.Args)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	Logger.info.Printf("Running command and waiting for it to finish...")
	tty, err := pty.Start(cmd)
	if err != nil {
		Logger.error.Fatal("Error start cmd", err)
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
		Logger.error.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	Logger.info.Printf("Done with command %s", cmd.Args)
	if err != nil {
		Logger.error.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		Logger.error.Fatal("failed to capture stdout or stderr\n")
	}

	outStr := string(stdoutBuf.Bytes())

	return outStr
}

func MatchCommandOrExt(toSearch []string, cmd string, ext string) bool {
	Logger.debug.Printf("MatchCommandOrExt got command: %s", strings.Join(toSearch, " "))
	if len(toSearch) == 0 {
		return false
	}
	file := filepath.Base(toSearch[0])
	return (file == cmd || filepath.Ext(file) == ext)
}

func Cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		Logger.error.Fatal("Error getting current working directory")
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

func readIgnoreFile(f string) []string {
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

func getNameOrEmpty(s INameable) string {
	n := ""
	if s != nil {
		n = s.Name() //Resource can be nil
	}
	return n
}

func runShell(cmd string, args []string) (*exec.Cmd, *os.File) {
	shell := exec.Command(cmd, args...)

	shell.Stderr = os.Stderr
	shell.Stdout = os.Stdout

	Logger.info.Printf("Running command and waiting for it to finish...")
	tty, err := pty.Start(shell)
	if err != nil {
		Logger.error.Fatal("Error start shell", err)
	}

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, tty); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	//p.oldState, err = terminal.MakeRaw(int(os.Stdin.Fd()))
	//if err != nil {
	//	panic(err)
	//}

	go func() {
		if _, err := io.Copy(os.Stdout, tty); err != nil {
			Logger.error.Fatal(err)
		}
	}()

	Logger.info.Printf("Done with command %s", shell.Args)
	return shell, tty
}

func runShellCmds(shell *exec.Cmd, tty *os.File, cb func([]string) []string) {
	go func() {
		stdin := bufio.NewReader(os.Stdin)
		for {
			line, _, err := stdin.ReadLine()
			if err != nil {
				if err == io.EOF {
					shell.Process.Kill()
					tty.Close()
				} else {
					Logger.error.Fatal(err)
				}
			}
			if cmd := cb(strings.Split(string(line), " ")); len(cmd) > 0 {
				tty.Write(append([]byte(strings.Join(cmd, " ")), '\n'))
			}
		}

	}()
	Logger.info.Printf("Waiting for shell to exit")
	shell.Wait()
	log.Printf("Shell exited")
}
