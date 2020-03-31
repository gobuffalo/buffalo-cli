package build

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
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
	fn := func(ctx context.Context, root string, cmd *exec.Cmd) error {
		act = cmd.Args
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
	fn := func(ctx context.Context, root string, args []string) error {
		return exp
	}

	p := buildtest.Builder(fn)
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			p,
		}
	})

	act := bc.Main(context.Background(), ".", []string{p.PluginName()})
	r.Equal(exp, act)
}

<<<<<<< HEAD
=======
func Test_Cmd_Main_ValidateTemplates(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}
	fn := func(root string) error {
		return nil
	}

	p := buildtest.TemplatesValidator(fn)
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			p,
		}
	})

	err := bc.Main(context.Background(), ".", nil)
	r.NoError(err)
}

func Test_Cmd_Main_ValidateTemplates_err(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	exp := fmt.Errorf("boom")
	fn := func(root string) error {
		return exp
	}

	p := buildtest.TemplatesValidator(fn)
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			p,
		}
	})

	err := bc.Main(context.Background(), ".", nil)
	r.Equal(exp, err)
}

>>>>>>> breaking the wall
func Test_Cmd_Main_BeforeBuilders(t *testing.T) {

	table := []struct {
		name string
		root string
		exp  []string
		err  error
	}{
		{name: "happy", root: ".", exp: []string{"-v"}},
		{name: "sad", root: ".", exp: []string{"-v"}, err: fmt.Errorf("boom")},
	}

	for _, tt := range table {
		t.Run(tt.name, func(st *testing.T) {
			r := require.New(st)

			var act []string
			fn := func(ctx context.Context, root string, args []string) error {
				act = args
				return tt.err
			}

			plugs := plugins.Plugins{buildtest.BeforeBuilder(fn)}

			bc := &Cmd{
				pluginsFn: func() []plugins.Plugin {
					return plugs
				},
			}

			err := bc.Main(context.Background(), tt.root, tt.exp)
			r.True(errors.Is(err, tt.err))
			r.Equal(tt.exp, act)
		})
	}

}

func Test_Cmd_Main_AfterBuilders(t *testing.T) {

	table := []struct {
		name string
		root string
		exp  []string
		err  error
	}{
		{name: "happy", root: ".", exp: []string{"-v"}},
		{name: "sad", root: ".", exp: []string{"-v"}, err: fmt.Errorf("boom")},
	}

	for _, tt := range table {
		t.Run(tt.name, func(st *testing.T) {
			r := require.New(st)

			var act []string
			fn := func(ctx context.Context, root string, args []string, err error) error {
				act = args
				return tt.err
			}

			plugs := plugins.Plugins{buildtest.AfterBuilder(fn)}

<<<<<<< HEAD
	b := &beforeBuilder{err: fmt.Errorf("science fiction twin")}
	a := &afterBuilder{}
	plugs := plugins.Plugins{a, b, &bladeRunner{}}
=======
			bc := &Cmd{
				pluginsFn: func() []plugins.Plugin {
					return plugs
				},
			}
>>>>>>> build stuff

			err := bc.Main(context.Background(), tt.root, tt.exp)
			r.Equal(tt.err, err)
			r.Equal(tt.exp, act)
		})
	}

<<<<<<< HEAD
	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.Error(err)
	r.Contains(err.Error(), b.err.Error())
=======
>>>>>>> build stuff
}
