package version

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	bufcli "github.com/gobuffalo/buffalo-cli/v2"
	"github.com/stretchr/testify/require"
)

func Test_Cmd(t *testing.T) {
	r := require.New(t)

	vc := &Cmd{}

	bb := &bytes.Buffer{}
	vc.SetStdout(bb)

	ctx := context.Background()

	args := []string{}

	err := vc.Main(ctx, ".", args)
	r.NoError(err)

	r.Contains(bb.String(), bufcli.Version)
}

func Test_Cmd_JSON(t *testing.T) {
	r := require.New(t)

	vc := &Cmd{}

	bb := &bytes.Buffer{}
	vc.SetStdout(bb)

	ctx := context.Background()

	args := []string{"--json"}

	err := vc.Main(ctx, ".", args)
	r.NoError(err)

	r.Contains(bb.String(), fmt.Sprintf("%q: %q", "version", bufcli.Version))
}
