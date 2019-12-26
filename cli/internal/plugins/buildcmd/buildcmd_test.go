package buildcmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here"
	"github.com/stretchr/testify/require"
)

type ref struct {
	here.Info
	t   *testing.T
	PWD string
}

func (r ref) Close() {
	os.RemoveAll(filepath.Join(r.Root, "main.build.go"))
	os.RemoveAll(filepath.Join(r.Root, "bin"))
	os.Chdir(r.PWD)
}

func newRef(t *testing.T, name string) ref {
	t.Helper()

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(pwd, "testdata", name)

	info := here.Info{
		Dir:        root,
		ImportPath: "github.com/markbates/coke",
		Name:       "main",
		Root:       root,
		Module: here.Module{
			Path:  "github.com/markbates/coke",
			Main:  true,
			Dir:   root,
			GoMod: filepath.Join(root, "go.mod"),
		},
	}

	r := ref{
		t:    t,
		Info: info,
		PWD:  pwd,
	}
	return r

}

func Test_BuildCmd_Subcommands(t *testing.T) {
	r := require.New(t)

	b := &builder{}
	all := plugins.Plugins{
		plugins.Background("foo"),
		&beforeBuilder{},
		b,
		&afterBuilder{},
		plugins.Background("bar"),
		&buildVersioner{},
		&templatesValidator{},
		&packager{},
	}

	bc := &BuildCmd{
		pluginsFn: all.ScopedPlugins,
	}

	plugs := bc.SubCommands()
	r.Len(plugs, 1)
	r.Equal(b, plugs[0])
}

func Test_BuildCmd_ScopedPlugins(t *testing.T) {
	r := require.New(t)

	all := plugins.Plugins{
		plugins.Background("foo"),
		&builder{},
		&beforeBuilder{},
		&afterBuilder{},
		plugins.Background("bar"),
		&buildVersioner{},
		&buildImporter{},
		&templatesValidator{},
		&packager{},
	}

	bc := &BuildCmd{
		pluginsFn: all.ScopedPlugins,
	}

	plugs := bc.ScopedPlugins()
	r.NotEqual(all, plugs)

	ep := plugins.Plugins(plugs).ExposedPlugins()

	tot := len(all) - 2
	r.Equal(tot, len(ep))

}
