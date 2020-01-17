package plugins

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
)

func Test_NewIO(t *testing.T) {
	ioe := NewIO()

	Equal(os.Stdin, ioe.Stdin())
	Equal(os.Stderr, ioe.Stderr())
	Equal(os.Stdout, ioe.Stdout())
}

func Test_WithStdin(t *testing.T) {
	Equal(os.Stdin, stdIO{}.Stdin())

	ctx := context.Background()
	ioe := CtxIO(ctx)
	Equal(os.Stdin, ioe.Stdin())

	in := &bytes.Buffer{}
	ctx = WithStdin(ctx, in)

	ioe = CtxIO(ctx)
	Equal(in, ioe.Stdin())
}

func Test_WithStdout(t *testing.T) {
	Equal(os.Stdout, stdIO{}.Stdout())

	ctx := context.Background()
	ioe := CtxIO(ctx)
	Equal(os.Stdout, ioe.Stdout())

	out := &bytes.Buffer{}
	ctx = WithStdout(ctx, out)

	ioe = CtxIO(ctx)
	Equal(out, ioe.Stdout())
}

func Test_WithStderr(t *testing.T) {
	Equal(os.Stderr, stdIO{}.Stderr())

	ctx := context.Background()
	ioe := CtxIO(ctx)
	Equal(os.Stderr, ioe.Stderr())

	out := &bytes.Buffer{}
	ctx = WithStderr(ctx, out)

	ioe = CtxIO(ctx)
	Equal(out, ioe.Stderr())
}

func Test_WithIO(t *testing.T) {
	ctx := context.Background()

	stdin := &bytes.Reader{}
	stdout := &strings.Builder{}
	stderr := &bytes.Buffer{}

	io := stdIO{
		Context: ctx,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
	}

	ctx = WithIO(ctx, io)

	ioe := CtxIO(ctx)
	Equal(stdin, ioe.Stdin())
	Equal(stdout, ioe.Stdout())
	Equal(stderr, ioe.Stderr())
}
