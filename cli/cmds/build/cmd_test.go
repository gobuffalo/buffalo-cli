package build

import (
	"context"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build/buildtest"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Subcommands(t *testing.T) {
	r := require.New(t)

	fn := func(ctx context.Context, root string, args []string) error {
		return nil
	}
	b := buildtest.Builder(fn)
	all := plugins.Plugins{
		background("foo"),
		buildtest.BeforeBuilder(nil),
		buildtest.Builder(nil),
		buildtest.AfterBuilder(nil),
		background("bar"),
		buildtest.Versioner(nil),
		buildtest.Packager(nil),
		buildtest.GoBuilder(nil),
	}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return all
		},
	}

	plugs := bc.SubCommands()
	r.Len(plugs, 1)
	r.Equal(b.PluginName(), plugs[0].PluginName())
}

func Test_Cmd_ScopedPlugins(t *testing.T) {
	r := require.New(t)

	all := plugins.Plugins{
		background("foo"),
		buildtest.Builder(nil),
		buildtest.BeforeBuilder(nil),
		buildtest.AfterBuilder(nil),
		background("bar"),
		buildtest.Versioner(nil),
		buildtest.Importer(nil),
		buildtest.GoBuilder(nil),
		buildtest.Packager(nil),
	}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return all
		},
	}

	plugs := bc.ScopedPlugins()
	r.NotEqual(all, plugs)

	ep := plugins.Plugins(plugs)

	tot := len(all) - 2
	r.Equal(tot, len(ep))

}
