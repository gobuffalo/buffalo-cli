package assets

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/here/there"
	"github.com/gobuffalo/meta/v2"
	"github.com/stretchr/testify/require"
)

func tempApp(t *testing.T, scripts map[string]string) here.Info {
	t.Helper()
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	sc := packageJSON{
		Scripts: scripts,
	}

	err = json.NewEncoder(io.MultiWriter(f, os.Stdout)).Encode(sc)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	return here.Info{
		Root: dir,
	}
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

func Test_Builder_Cmd_PackageJSON(t *testing.T) {
	r := require.New(t)

	info := tempApp(t, map[string]string{
		"build": "echo hi",
	})

	defer os.RemoveAll(info.Root)

	bc := &Builder{}
	ctx := context.Background()
	args := []string{}

	app, err := meta.New(info)
	r.NoError(err)

	c, err := bc.Cmd(app, ctx, args)
	r.NoError(err)

	r.Equal([]string{"npm", "run", "build"}, c.Args)
}

func Test_Builder_Cmd_PackageJSON_Yarn(t *testing.T) {
	r := require.New(t)

	info := tempApp(t, map[string]string{
		"build": "echo hi",
	})

	defer os.RemoveAll(info.Root)

	bc := &Builder{}
	ctx := context.Background()
	args := []string{}

	app, err := meta.New(info)
	r.NoError(err)

	app.With["yarn"] = true
	c, err := bc.Cmd(app, ctx, args)
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

	app, err := meta.New(info)
	r.NoError(err)
	c, err := bc.Cmd(app, ctx, args)
	r.NoError(err)

	r.Equal([]string{bc.webpackBin()}, c.Args)
}
