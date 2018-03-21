package shell

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

func NewStdinSwitch(stdin io.ReadWriter, stdout io.Writer) *StdinSwitch {
	r, w := io.Pipe()
	a := &StdinSwitch{
		r:      r,
		w:      w,
		stdin:  stdin,
		stdout: stdout,
		prompt: false,
	}
	a.Start()
	return a
}

type StdinSwitch struct {
	r      *io.PipeReader
	w      *io.PipeWriter
	stdin  io.Reader
	stdout io.Writer
	prompt bool
}

func (s *StdinSwitch) Start() {
	// Uses r, w to direct input and prevent issues where we can't prompt the user because
	// we are doing a a read in a background process
	go func() {
		for {
			var out []byte
			_, err := s.stdin.Read(out)
			if err != nil {
				log.Fatal(err)
			}
			if s.prompt {
				// We are currently in s.Prompt() and input shouldn't be forwarded
				continue
			}
			s.w.Write(out)
		}
	}()
}

func (s *StdinSwitch) Read(b []byte) (n int, err error) {
	n, err = s.r.Read(b)
	return n, err
}

func (s *StdinSwitch) Prompt(p string) string {
	s.prompt = true
	fmt.Fprint(s.stdout, p)
	out, err := bufio.NewReader(s.stdin).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	s.prompt = false
	return string(out)
}
