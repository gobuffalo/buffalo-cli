package build

import (
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/here"
	"github.com/stretchr/testify/require"
)

func newRef(t *testing.T, root string) here.Info {
	t.Helper()

	info := here.Info{
		Dir:        root,
		ImportPath: "github.com/markbates/coke",
		Name:       "coke",
		Root:       root,
		Module: here.Module{
			Path:  "github.com/markbates/coke",
			Main:  true,
			Dir:   root,
			GoMod: filepath.Join(root, "go.mod"),
		},
	}

	return info
}

func Test_Cmd_Subcommands(t *testing.T) {
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
		&bladeRunner{},
	}

	bc := &Cmd{
		pluginsFn: all.ScopedPlugins,
	}

	plugs := bc.SubCommands()
	r.Len(plugs, 1)
	r.Equal(b, plugs[0])
}

func Test_Cmd_ScopedPlugins(t *testing.T) {
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
		&bladeRunner{},
		&packager{},
	}

	bc := &Cmd{
		pluginsFn: all.ScopedPlugins,
	}

	plugs := bc.ScopedPlugins()
	r.NotEqual(all, plugs)

	ep := plugins.Plugins(plugs).ExposedPlugins()

	tot := len(all) - 2
	r.Equal(tot, len(ep))

}
