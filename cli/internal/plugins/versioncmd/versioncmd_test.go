package versioncmd

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/stretchr/testify/require"
)

func Test_VersionCmd(t *testing.T) {
	r := require.New(t)

	vc := &VersionCmd{}

	bb := &bytes.Buffer{}

	ctx := context.Background()
	ctx = plugins.WithStdout(ctx, bb)

	args := []string{}

	err := vc.Main(ctx, args)
	r.NoError(err)

	r.Contains(bb.String(), bufcli.Version)
}

func Test_VersionCmd_JSON(t *testing.T) {
	r := require.New(t)

	vc := &VersionCmd{}

	bb := &bytes.Buffer{}

	ctx := context.Background()
	ctx = plugins.WithStdout(ctx, bb)

	args := []string{"--json"}

	err := vc.Main(ctx, args)
	r.NoError(err)

	r.Contains(bb.String(), fmt.Sprintf("%q: %q", "version", bufcli.Version))
}
