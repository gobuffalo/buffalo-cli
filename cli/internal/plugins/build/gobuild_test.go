package build

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cmd_GoCmd(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx)
	r.NoError(err)

	cli := filepath.Join("bin", "buffalo-cli")
	if runtime.GOOS == "windows" {
		cli += ".exe"
	}
	exp := []string{"go", "build", "-o", cli}
	r.Equal(exp, cmd.Args)
}

func Test_Cmd_GoCmd_Bin(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{
		Bin: "cli",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx)
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
		Bin: "cli",
		Mod: "vendor",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx)
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
		Bin:  "cli",
		Tags: "a b c",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx)
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
		Bin:    "cli",
		Static: true,
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx)
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
		Bin:     "cli",
		LDFlags: "linky",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx)
	r.NoError(err)

	n := "cli"
	if runtime.GOOS == "windows" {
		n = "cli.exe"
	}

	exp := []string{"go", "build", "-o", n, "-ldflags", "linky"}
	r.Equal(exp, cmd.Args)
}
