package shell

import (
	"bufio"
	"bytes"
	"io"
)

func NewAsk(i io.Reader, o io.Writer) io.ReadWriter {
	a := &Ask{
		input:  i,
		output: o,
	}
	return a
}

type Ask struct {
	input  io.Reader
	output io.Writer
}

func (a *Ask) Input() *bufio.Reader {
	return bufio.NewReader(new(bytes.Buffer))
}

func (a *Ask) Write(b []byte) (int, error) {
	return a.output.Write(b)
}

func (a *Ask) Read(b []byte) (int, error) {
	n, err := a.input.Read(b)
	if err != nil {
		return n, err
	}
	return n, nil
}
