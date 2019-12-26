package buildcmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/stretchr/testify/require"
)

func Test_BuildCmd_Main(t *testing.T) {
	r := require.New(t)

	info := newRef(t, "")

	bc := &BuildCmd{}
	bc.WithHereInfo(info)

	exp := []string{"go", "build", "-o", "bin/coke"}

	br := &bladeRunner{}
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{br}
	})

	var args []string
	err := bc.Main(context.Background(), args)
	r.NoError(err)
	r.NotNil(br.cmd)
	r.Equal(exp, br.cmd.Args)
}

func Test_BuildCmd_Main_SubCommand(t *testing.T) {
	r := require.New(t)

	p := &builder{name: "foo"}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &BuildCmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	args := []string{p.name, "a", "b", "c"}

	err := bc.Main(context.Background(), args)
	r.NoError(err)
	r.Equal([]string{"a", "b", "c"}, p.args)
}

func Test_BuildCmd_Main_SubCommand_err(t *testing.T) {
	r := require.New(t)

	p := &builder{name: "foo", err: fmt.Errorf("error")}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &BuildCmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	args := []string{p.name, "a", "b", "c"}

	err := bc.Main(context.Background(), args)
	r.Error(err)
}

func Test_BuildCmd_Main_ValidateTemplates(t *testing.T) {
	r := require.New(t)

	p := &templatesValidator{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	info := newRef(t, "")

	bc := &BuildCmd{
		Info:      info,
		pluginsFn: plugs.ScopedPlugins,
	}

	args := []string{}

	err := bc.Main(context.Background(), args)
	r.NoError(err)
	r.Equal(bc.Info.Root, p.root)
}

func Test_BuildCmd_Main_ValidateTemplates_err(t *testing.T) {
	r := require.New(t)

	p := &templatesValidator{err: fmt.Errorf("error")}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &BuildCmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	args := []string{}

	err := bc.Main(context.Background(), args)
	r.Error(err)
}

func Test_BuildCmd_Main_BeforeBuilders(t *testing.T) {
	r := require.New(t)

	p := &beforeBuilder{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &BuildCmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), args)
	r.NoError(err)
}

func Test_BuildCmd_Main_BeforeBuilders_err(t *testing.T) {
	r := require.New(t)

	p := &beforeBuilder{err: fmt.Errorf("error")}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &BuildCmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), args)
	r.Error(err)
}

func Test_BuildCmd_Main_AfterBuilders(t *testing.T) {
	r := require.New(t)

	p := &afterBuilder{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &BuildCmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), args)
	r.NoError(err)
}

func Test_BuildCmd_Main_AfterBuilders_err(t *testing.T) {
	r := require.New(t)

	b := &beforeBuilder{err: fmt.Errorf("error")}
	a := &afterBuilder{}
	plugs := plugins.Plugins{a, b, &bladeRunner{}}

	bc := &BuildCmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), args)
	r.Error(err)
	r.Equal(err, a.err)
}
