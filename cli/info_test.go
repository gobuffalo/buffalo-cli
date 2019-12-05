package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/stretchr/testify/require"
)

type snow struct{}

func (s snow) Name() string {
	return "snow"
}

func (s snow) Info(ctx context.Context, args []string) error {
	out := cmdx.Stdout(ctx)
	out.Write([]byte("informer"))
	return nil
}

func Test_Buffalo_Info(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}

	ctx := context.Background()
	ctx = cmdx.WithStdout(ctx, bb)

	buffalo, err := New(ctx)
	r.NoError(err)

	buffalo.Plugins = plugins.Plugins{
		snow{},
	}

	err = buffalo.Info(ctx, []string{})
	r.NoError(err)

	out := bb.String()
	r.Contains(out, "Buffalo (CLI)")
	r.Contains(out, "informer")
}
