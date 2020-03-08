package build

import (
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build/buildtest"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Subcommands(t *testing.T) {
	r := require.New(t)

	b := buildtest.Builder(nil)
	all := plugins.Plugins{
		background("foo"),
		buildtest.BeforeBuilder(nil),
		buildtest.Builder(nil),
		buildtest.AfterBuilder(nil),
		background("bar"),
		buildtest.Versioner(nil),
		buildtest.TemplatesValidator(nil),
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
	r.Equal(b, plugs[0])
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
		buildtest.TemplatesValidator(nil),
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
