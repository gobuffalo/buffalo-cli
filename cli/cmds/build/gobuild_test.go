package build

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build/buildtest"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_GoCmd(t *testing.T) {
	r := require.New(t)

	var act []string
	fn := func(ctx context.Context, root string, args []string) error {
		act = args
		return nil
	}

	bc := &Cmd{}
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			buildtest.GoBuilder(fn),
		}
	})

	ctx := context.Background()
	err := bc.Main(ctx, ".", nil)
	r.NoError(err)

	cli := filepath.Join("bin", "build")
	if runtime.GOOS == "windows" {
		cli += ".exe"
	}
	exp := []string{"go", "build", "-o", cli}
	r.Equal(exp, args)
}

func Test_Cmd_GoCmd_Bin(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{
		bin: "cli",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)

	n := "cli"
	if runtime.GOOS == "windows" {
		n = "cli.exe"
	}

	exp := []string{"go", "build", "-o", n}
	r.Equal(exp, cmd.Args)
}

func Test_Cmd_GoCmd_Mod(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{
		bin: "cli",
		mod: "vendor",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)

	n := "cli"
	if runtime.GOOS == "windows" {
		n = "cli.exe"
	}

	exp := []string{"go", "build", "-o", n, "-mod", "vendor"}
	r.Equal(exp, cmd.Args)
}

func Test_Cmd_GoCmd_Tags(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{
		bin:  "cli",
		tags: "a b c",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)

	n := "cli"
	if runtime.GOOS == "windows" {
		n = "cli.exe"
	}

	exp := []string{"go", "build", "-o", n, "-tags", "a b c"}
	r.Equal(exp, cmd.Args)
}

func Test_Cmd_GoCmd_Static(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{
		bin:    "cli",
		static: true,
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)

	n := "cli"
	if runtime.GOOS == "windows" {
		n = "cli.exe"
	}

	exp := []string{"go", "build", "-o", n, "-ldflags", "-linkmode external -extldflags \"-static\""}
	r.Equal(exp, cmd.Args)
}

func Test_Cmd_GoCmd_LDFlags(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{
		bin:     "cli",
		ldFlags: "linky",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)

	n := "cli"
	if runtime.GOOS == "windows" {
		n = "cli.exe"
	}

	exp := []string{"go", "build", "-o", n, "-ldflags", "linky"}
	r.Equal(exp, cmd.Args)
}
