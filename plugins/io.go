package plugins

import (
	"io"
	"os"
)

type IO interface {
	StderrGetter
	StderrSetter
	StdinGetter
	StdinSetter
	StdoutGetter
	StdoutSetter
}

type IOGetters interface {
	StderrGetter
	StdinGetter
	StdoutGetter
}

type IOSetters interface {
	StderrSetter
	StdinSetter
	StdoutSetter
}

type StdinGetter interface {
	Stdin() io.Reader
}

type StdoutGetter interface {
	Stdout() io.Writer
}

type StderrGetter interface {
	Stderr() io.Writer
}

type StdinSetter interface {
	SetStdin(r io.Reader)
}

type StdoutSetter interface {
	SetStdout(w io.Writer)
}

type StderrSetter interface {
	SetStderr(w io.Writer)
}

func NewIO() IO {
	return &stdIO{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

type stdIO struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (i *stdIO) Stdin() io.Reader {
	if i.stdin == nil {
		return os.Stdin
	}
	return i.stdin
}

func (i *stdIO) Stdout() io.Writer {
	if i.stdout == nil {
		return os.Stdout
	}
	return i.stdout
}

func (i *stdIO) Stderr() io.Writer {
	if i.stderr == nil {
		return os.Stderr
	}
	return i.stderr
}

func (i *stdIO) SetStdin(r io.Reader) {
	i.stdin = r
}

func (i *stdIO) SetStdout(w io.Writer) {
	i.stdout = w
}

func (i *stdIO) SetStderr(w io.Writer) {
	i.stderr = w
}

func SetIO(in IO, p Plugin) {
	if stdin, ok := p.(StdinSetter); ok {
		stdin.SetStdin(in.Stdin())
	}
	if stdout, ok := p.(StdoutSetter); ok {
		stdout.SetStdout(in.Stdout())
	}
	if stderr, ok := p.(StderrSetter); ok {
		stderr.SetStderr(in.Stderr())
	}
}
