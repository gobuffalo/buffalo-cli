package build

import (
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build/buildtest"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_GoBuilder(t *testing.T) {
	cli := func(n ...string) string {
		cli := filepath.Join(n...)
		if runtime.GOOS == "windows" {
			cli += ".exe"
		}
		return cli
	}

	mainFolder := filepath.Join("cmd", "build")
	if runtime.GOOS == "windows" {
		mainFolder = "\\" + mainFolder
	} else {
		mainFolder = "./" + mainFolder
	}

	table := []struct {
		args []string
		exp  []string
	}{
		{
			exp: []string{"go", "build", "-o", cli("bin", "build"), mainFolder},
		},
		{
			args: []string{"-o", filepath.Join("bin", "foo")},
			exp:  []string{"go", "build", "-o", cli("bin", "foo"), mainFolder},
		},
		{
			args: []string{"--mod", "vendor"},
			exp:  []string{"go", "build", "-o", cli("bin", "build"), "-mod", "vendor", mainFolder},
		},
		{
			args: []string{"--tags", "a b c"},
			exp:  []string{"go", "build", "-o", cli("bin", "build"), "-tags", "a b c", mainFolder},
		},
		{
			args: []string{"--static"},
			exp:  []string{"go", "build", "-o", cli("bin", "build"), "-ldflags", "-linkmode external -extldflags \"-static\"", mainFolder},
		},
		{
			args: []string{"--ldflags", "linky"},
			exp:  []string{"go", "build", "-o", cli("bin", "build"), "-ldflags", "linky", mainFolder},
		},
	}

	for _, tt := range table {
		t.Run(strings.Join(tt.args, " "), func(st *testing.T) {
			r := require.New(st)

			bc := &Cmd{}

			var act []string
			fn := func(ctx context.Context, root string, cmd *exec.Cmd) error {
				act = cmd.Args
				return nil
			}
			bc.WithPlugins(func() []plugins.Plugin {
				return []plugins.Plugin{
					buildtest.GoBuilder(fn),
				}
			})

			ctx := context.Background()
			err := bc.Main(ctx, ".", tt.args)
			r.NoError(err)

			r.Equal(tt.exp, act)
		})
	}
}
