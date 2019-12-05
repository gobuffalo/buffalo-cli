package plugins

import (
	"bytes"
	"context"
	"testing"

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

func Test_Info(t *testing.T) {
	r := require.New(t)

	plugs := Plugins{
		snow{},
	}

	bb := &bytes.Buffer{}

	ctx := context.Background()
	ctx = cmdx.WithStdout(ctx, bb)

	err := plugs.Info(ctx, []string{})
	r.NoError(err)

	r.Contains(bb.String(), "informer")
}
