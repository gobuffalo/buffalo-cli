package plugins

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewIO(t *testing.T) {
	r := require.New(t)

	ioe := NewIO()

	r.Equal(os.Stdin, ioe.Stdin())
	r.Equal(os.Stderr, ioe.Stderr())
	r.Equal(os.Stdout, ioe.Stdout())
}

func Test_WithStdin(t *testing.T) {
	r := require.New(t)

	r.Equal(os.Stdin, stdIO{}.Stdin())

	ctx := context.Background()
	ioe := CtxIO(ctx)
	r.Equal(os.Stdin, ioe.Stdin())

	in := &bytes.Buffer{}
	ctx = WithStdin(ctx, in)

	ioe = CtxIO(ctx)
	r.Equal(in, ioe.Stdin())
}

func Test_WithStdout(t *testing.T) {
	r := require.New(t)

	r.Equal(os.Stdout, stdIO{}.Stdout())

	ctx := context.Background()
	ioe := CtxIO(ctx)
	r.Equal(os.Stdout, ioe.Stdout())

	out := &bytes.Buffer{}
	ctx = WithStdout(ctx, out)

	ioe = CtxIO(ctx)
	r.Equal(out, ioe.Stdout())
}

func Test_WithStderr(t *testing.T) {
	r := require.New(t)

	r.Equal(os.Stderr, stdIO{}.Stderr())

	ctx := context.Background()
	ioe := CtxIO(ctx)
	r.Equal(os.Stderr, ioe.Stderr())

	out := &bytes.Buffer{}
	ctx = WithStderr(ctx, out)

	ioe = CtxIO(ctx)
	r.Equal(out, ioe.Stderr())
}

func Test_WithIO(t *testing.T) {
	r := require.New(t)

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
	r.Equal(stdin, ioe.Stdin())
	r.Equal(stdout, ioe.Stdout())
	r.Equal(stderr, ioe.Stderr())
}
