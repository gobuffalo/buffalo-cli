package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/stretchr/testify/require"
)

func Test_Buffalo_Version(t *testing.T) {
	r := require.New(t)

	ctx := context.Background()
	bb := &bytes.Buffer{}
	ctx = cmdx.WithStdout(ctx, bb)

	buffalo, err := New(ctx)
	r.NoError(err)

	err = buffalo.Main(ctx, []string{"version"})
	r.NoError(err)

	out := strings.TrimSpace(bb.String())
	r.Equal(out, bufcli.Version)
}

func Test_Buffalo_Version_JSON(t *testing.T) {
	r := require.New(t)

	ctx := context.Background()
	bb := &bytes.Buffer{}
	ctx = cmdx.WithStdout(ctx, bb)

	buffalo, err := New(ctx)
	r.NoError(err)

	err = buffalo.Main(ctx, []string{"version", "--json"})
	r.NoError(err)

	out := strings.TrimSpace(bb.String())
	r.Contains(out, "version")
	r.Contains(out, bufcli.Version)

}
