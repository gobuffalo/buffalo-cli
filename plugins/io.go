package plugins

import (
	"context"
	"io"
)

// CtxIO returns a working IO implmentation which
// defaults to using os.Stdin, os.Stdout, and os.Stderr.
// If the context itself implements IO, the it is returned.
// Next, if the context contains an "io" value that implements
// IO, then that value is returned.
// If not IO implementation is found, a default implementation
// using the standard IO is returned.
func CtxIO(ctx context.Context) IO {
	if i, ok := ctx.(IO); ok {
		return i
	}
	if i, ok := ctx.Value("io").(IO); ok {
		return i
	}
	return NewIO()
}

// WithIO wraps the given context and IO with a new context
// that also implements the given IO.
func WithIO(ctx context.Context, i IO) context.Context {
	return stdIO{
		Context: ctx,
		stdout:  i.Stdout(),
		stderr:  i.Stderr(),
		stdin:   i.Stdin(),
	}
}

// WithStdin returns a new context that implements IO with
// the given io.Reader representing Stdin.
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
