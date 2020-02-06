package plugins

import (
	"context"
	"io"
	"os"
)

func NewIO() IO {
	return stdIO{
		Context: context.Background(),
		stdin:   os.Stdin,
		stdout:  os.Stdout,
		stderr:  os.Stderr,
	}
}

type stdIO struct {
	context.Context
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (i stdIO) Stdin() io.Reader {
	if i.stdin != nil {
		return i.stdin
	}
	if x, ok := i.Context.(StdinGetter); ok {
		return x.Stdin()
	}
	return os.Stdin
}

func (i stdIO) Stdout() io.Writer {
	if i.stdout != nil {
		return i.stdout
	}
	if x, ok := i.Context.(StdoutGetter); ok {
		return x.Stdout()
	}
	return os.Stdout
}

func (i stdIO) Stderr() io.Writer {
	if i.stderr != nil {
		return i.stderr
	}
	if x, ok := i.Context.(StderrGetter); ok {
		return x.Stderr()
	}
	return os.Stderr
}
