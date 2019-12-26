package plugins

import (
	"context"
	"io"
	"os"
)

type IO interface {
	StderrGetter
	StdinGetter
	StdoutGetter
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

func CtxIO(ctx context.Context) IO {
	if i, ok := ctx.(IO); ok {
		return i
	}
	if i, ok := ctx.Value("io").(IO); ok {
		return i
	}
	return NewIO()
}

func WithIO(ctx context.Context, i IO) context.Context {
	return stdIO{
		Context: ctx,
		stdout:  i.Stdout(),
		stderr:  i.Stderr(),
		stdin:   i.Stdin(),
	}
}

func WithStdin(ctx context.Context, stdin io.Reader) context.Context {
	i := CtxIO(ctx)
	return stdIO{
		Context: ctx,
		stdout:  i.Stdout(),
		stderr:  i.Stderr(),
		stdin:   stdin,
	}
}

func WithStdout(ctx context.Context, stdout io.Writer) context.Context {
	i := CtxIO(ctx)
	return stdIO{
		Context: ctx,
		stdout:  stdout,
		stderr:  i.Stderr(),
		stdin:   i.Stdin(),
	}
}

func WithStderr(ctx context.Context, stderr io.Writer) context.Context {
	i := CtxIO(ctx)
	return stdIO{
		Context: ctx,
		stdout:  i.Stdout(),
		stderr:  stderr,
		stdin:   i.Stdin(),
	}
}
