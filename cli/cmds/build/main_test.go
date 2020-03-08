package build

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build/buildtest"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Main(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	bn := filepath.Join("bin", "build")
	if runtime.GOOS == "windows" {
		bn += ".exe"
	}
	exp := []string{"go", "build", "-o", bn}

	var act []string
	fn := func(ctx context.Context, root string, args []string) error {
		act = args
		return nil
	}
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			buildtest.GoBuilder(fn),
		}
	})

	var args []string
	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
	r.Equal(exp, act)
}

func Test_Cmd_Main_SubCommand(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	var act []string
	fn := func(ctx context.Context, root string, args []string) error {
		act = args
		return nil
	}

	p := buildtest.Builder(fn)
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			p,
		}
	})

	args := []string{p.PluginName(), "a", "b", "c"}

	err := bc.Main(context.Background(), "", args)
	r.NoError(err)
	r.Equal(args[1:], act)
}

func Test_Cmd_Main_SubCommand_err(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	exp := fmt.Errorf("boom")
	var act []string
	fn := func(ctx context.Context, root string, args []string) error {
		act = args
		return exp
	}

	p := buildtest.Builder(fn)
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			p,
		}
	})

	args := []string{p.PluginName(), "a", "b", "c"}

	err := bc.Main(context.Background(), "", args)
	r.Error(err)
	r.Equal(exp, err)
}

<<<<<<< HEAD
=======
func Test_Cmd_Main_ValidateTemplates(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	var act string
	fn := func(root string) error {
		act = root
		return nil
	}

	p := buildtest.TemplatesValidator(fn)
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			p,
		}
	})

	err := bc.Main(context.Background(), "", nil)
	r.NoError(err)
}

func Test_Cmd_Main_ValidateTemplates_err(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	exp := fmt.Errorf("boom")
	var act string
	fn := func(root string) error {
		act = root
		return exp
	}

	p := buildtest.TemplatesValidator(fn)
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			p,
		}
	})

	err := bc.Main(context.Background(), "", nil)
	r.Error(err)
	r.Equal(exp, err)
}

>>>>>>> breaking the wall
func Test_Cmd_Main_BeforeBuilders(t *testing.T) {
	r := require.New(t)

	p := &beforeBuilder{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
}

func Test_Cmd_Main_BeforeBuilders_err(t *testing.T) {
	r := require.New(t)

	p := &beforeBuilder{err: fmt.Errorf("error")}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.Error(err)
}

func Test_Cmd_Main_AfterBuilders(t *testing.T) {
	r := require.New(t)

	p := &afterBuilder{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
}

func Test_Cmd_Main_AfterBuilders_err(t *testing.T) {
	r := require.New(t)

	b := &beforeBuilder{err: fmt.Errorf("science fiction twin")}
	a := &afterBuilder{}
	plugs := plugins.Plugins{a, b, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.Error(err)
	r.Contains(err.Error(), b.err.Error())
}
