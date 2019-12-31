package testcmd

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/stretchr/testify/require"
)

func cmpCmd(t *testing.T, exp string, cmd *exec.Cmd) {
	act := strings.Join(cmd.Args, " ")
	if act != exp {
		t.Fatalf("expected '%s', got '%s'", exp, act)
	}
}

func Test_TestCmd_Cmd(t *testing.T) {
	r := require.New(t)

	ctx := context.Background()
	args := []string{}

	tc := &TestCmd{}

	cmd, err := tc.Cmd(ctx, args)
	r.NoError(err)

	cmpCmd(t, "go test ./...", cmd)
}

func Test_TestCmd_Cmd_buildArgs(t *testing.T) {
	r := require.New(t)

	ctx := context.Background()
	args := []string{"-tags", "you're"}

	tc := &TestCmd{}

	cmd, err := tc.Cmd(ctx, args)
	r.NoError(err)

	cmpCmd(t, "go test -tags you're", cmd)

	args = append(args, "-run", "Foo", "-tags", "it", "-v", "./...")

	cmd, err = tc.Cmd(ctx, args)
	r.NoError(err)

	cmpCmd(t, "go test -tags you're it -run Foo -v ./...", cmd)

	tc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			&tagger{
				tags: []string{"-tags", "i", "-failfast"},
			},
		}
	})

	cmd, err = tc.Cmd(ctx, args)
	r.NoError(err)
	cmpCmd(t, "go test -tags i you're it -failfast -run Foo -v ./...", cmd)
}
