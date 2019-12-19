package assets

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here/there"
	"github.com/stretchr/testify/require"
)

func Test_Builder_Build_Help(t *testing.T) {
	r := require.New(t)

	bc := &Builder{}
	ctx := context.Background()

	stdout := &bytes.Buffer{}
	ctx = plugins.WithStdout(ctx, stdout)
	args := []string{"-h"}

	err := bc.Build(ctx, args)
	r.NoError(err)
	r.Contains(stdout.String(), bc.Description())
}

func Test_Builder_Build(t *testing.T) {
	r := require.New(t)

	info := tempApp(t, map[string]string{
		"build": "echo wolverine",
	})

	defer os.RemoveAll(info.Root)

	pwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(pwd)

	os.Chdir(info.Root)

	bc := &Builder{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "here.Current", info)

	stdout := &bytes.Buffer{}
	ctx = plugins.WithStdout(ctx, stdout)

	args := []string{}

	err = bc.Build(ctx, args)
	r.NoError(err)
	r.Contains(stdout.String(), "wolverine")
}

func Test_Builder_Build_Skip(t *testing.T) {
	r := require.New(t)

	info := tempApp(t, map[string]string{
		"build": "echo wolverine",
	})

	defer os.RemoveAll(info.Root)

	pwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(pwd)

	os.Chdir(info.Root)

	bc := &Builder{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "here.Current", info)

	stdout := &bytes.Buffer{}
	ctx = plugins.WithStdout(ctx, stdout)

	args := []string{"--skip-assets"}

	err = bc.Build(ctx, args)
	r.NoError(err)
	r.Empty(stdout.String())
}

func Test_Builder_Cmd_PackageJSON(t *testing.T) {
	r := require.New(t)

	info := tempApp(t, map[string]string{
		"build": "echo hi",
	})

	defer os.RemoveAll(info.Root)

	bc := &Builder{}
	ctx := context.Background()
	args := []string{}

	c, err := bc.Cmd(info.Dir, ctx, args)
	r.NoError(err)

	r.Equal([]string{"npm", "run", "build"}, c.Args)
}

func Test_Builder_Cmd_PackageJSON_Yarn(t *testing.T) {
	r := require.New(t)

	info := tempApp(t, map[string]string{
		"build": "echo hi",
	})

	defer os.RemoveAll(info.Root)

	bc := &Builder{
		Tool: "yarnpkg",
	}
	ctx := context.Background()
	args := []string{}

	c, err := bc.Cmd(info.Root, ctx, args)
	r.NoError(err)

	r.Equal([]string{"yarnpkg", "run", "build"}, c.Args)
}

func Test_Builder_Cmd_Webpack_Fallthrough(t *testing.T) {
	r := require.New(t)

	bc := &Builder{}

	ctx := context.Background()
	args := []string{}

	info, err := there.Current()
	r.NoError(err)

	c, err := bc.Cmd(info.Root, ctx, args)
	r.NoError(err)

	r.Equal([]string{bc.webpackBin()}, c.Args)
}
