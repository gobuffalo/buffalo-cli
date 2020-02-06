package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Main(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	exp := []string{"go", "test", "./..."}

	br := &bladeRunner{}
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{br}
	})

	var args []string
	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
	r.NotNil(br.cmd)
	r.Equal(exp, br.cmd.Args)
}

func Test_Cmd_Main_SubCommand(t *testing.T) {
	r := require.New(t)

	p := &tester{name: "foo"}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	args := []string{p.name, "a", "b", "c"}

	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
	r.Equal([]string{"a", "b", "c"}, p.args)
}

func Test_Cmd_Main_SubCommand_err(t *testing.T) {
	r := require.New(t)

	p := &tester{name: "foo", err: fmt.Errorf("error")}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	args := []string{p.name, "a", "b", "c"}

	err := bc.Main(context.Background(), ".", args)
	r.Error(err)
}

func Test_Cmd_Main_BeforeTesters(t *testing.T) {
	r := require.New(t)

	p := &beforeTester{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
}

func Test_Cmd_Main_BeforeTesters_err(t *testing.T) {
	r := require.New(t)

	p := &beforeTester{err: fmt.Errorf("error")}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.Error(err)
}

func Test_Cmd_Main_AfterTesters(t *testing.T) {
	r := require.New(t)

	p := &afterTester{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
}

func Test_Cmd_Main_AfterTesters_err(t *testing.T) {
	r := require.New(t)

	b := &beforeTester{err: fmt.Errorf("error")}
	a := &afterTester{}
	plugs := plugins.Plugins{a, b, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.Error(err)
	r.Equal(err, a.err)
}
