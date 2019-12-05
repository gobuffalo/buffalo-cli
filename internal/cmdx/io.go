package cmdx

import (
	"context"
	"io"
	"os"
)

const (
	stdin  = "stdin"
	stdout = "stdout"
	stderr = "stderr"
)

func WithStdin(ctx context.Context, r io.Reader) context.Context {
	return context.WithValue(ctx, stdin, r)
}

func Stdin(ctx context.Context) io.Reader {
	if r, ok := ctx.Value(stdin).(io.Reader); ok {
		return r
	}
	return os.Stdin
}

func WithStdout(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, stdout, w)
}

func Stdout(ctx context.Context) io.Writer {
	if w, ok := ctx.Value(stdout).(io.Writer); ok {
		return w
	}
	return os.Stdout
}

func WithStderr(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, stderr, w)
}

func Stderr(ctx context.Context) io.Writer {
	if w, ok := ctx.Value(stderr).(io.Writer); ok {
		return w
	}
	return os.Stderr
}
