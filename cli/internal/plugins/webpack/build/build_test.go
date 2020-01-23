package build

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack/internal/scripts"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/here"
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

	err = json.NewEncoder(f).Encode(sc)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	return here.Info{
		Root: dir,
		Dir:  dir,
	}
}

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

	defer os.RemoveAll(info.Dir)

	pwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(pwd)

	os.Chdir(info.Dir)

	bc := &Builder{}
	bc.WithHereInfo(info)

	br := &bladeRunner{}
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{br}
	})

	exp := []string{"npm", "run", "build"}

	ctx := context.Background()

	args := []string{}

	err = bc.Build(ctx, args)
	r.NoError(err)
	r.Equal(exp, br.cmd.Args)
}

func Test_Builder_Build_Skip(t *testing.T) {
	r := require.New(t)

	bc := &Builder{}
	br := &bladeRunner{}
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			br,
			&tooler{tool: "npm"},
		}
	})

	ctx := context.Background()

	stdout := &bytes.Buffer{}
	ctx = plugins.WithStdout(ctx, stdout)

	args := []string{"--skip-webpack"}

	err := bc.Build(ctx, args)
	r.NoError(err)
	r.Empty(stdout.String())
}

func Test_Builder_Cmd_PackageJSON(t *testing.T) {
	r := require.New(t)

	info := tempApp(t, map[string]string{
		"build": "echo hi",
	})

	defer os.RemoveAll(info.Dir)

	bc := &Builder{}
	ctx := context.Background()
	args := []string{}

	c, err := bc.cmd(ctx, info.Dir, args)
	r.NoError(err)

	r.Equal([]string{"npm", "run", "build"}, c.Args)
}

func Test_Builder_Cmd_PackageJSON_Yarn(t *testing.T) {
	r := require.New(t)

	info := tempApp(t, map[string]string{
		"build": "echo hi",
	})

	defer os.RemoveAll(info.Dir)

	bc := &Builder{
		Tool: "yarnpkg",
	}
	ctx := context.Background()
	args := []string{}

	c, err := bc.cmd(ctx, info.Dir, args)
	r.NoError(err)

	r.Equal([]string{"yarnpkg", "run", "build"}, c.Args)
}

func Test_Builder_Cmd_Webpack_Fallthrough(t *testing.T) {
	r := require.New(t)

	bc := &Builder{}
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			&tooler{tool: "npm"},
		}
	})

	ctx := context.Background()
	args := []string{}

	info, err := here.Current()
	r.NoError(err)

	c, err := bc.cmd(ctx, info.Dir, args)
	r.NoError(err)

	r.Equal([]string{scripts.WebpackBin(info.Dir)}, c.Args)
}
