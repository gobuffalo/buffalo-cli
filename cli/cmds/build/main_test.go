package build

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build/buildtest"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

var (
	mainFileBeforeBuilder = func(ctx context.Context, root string, args []string) error {
		err := os.MkdirAll(filepath.Join(root, "cmd", "build"), 0777)
		if err != nil {
			return err
		}

		main := []byte(`package main
			func main() {
				
			}
		`)

		err = ioutil.WriteFile(filepath.Join(root, "cmd", "build", "main.go"), main, 0777)
		return err
	}

	cleanupAfterBuilder = func(ctx context.Context, root string, args []string, oerr error) error {
		err := os.RemoveAll(filepath.Join(root, "cmd"))
		if err != nil {
			return err
		}

		return nil
	}
)

func Test_Cmd_Main(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}
	bn := filepath.Join("bin", "build")

	if runtime.GOOS == "windows" {
		bn += ".exe"
	}

	mainFolder := filepath.Join("cmd", "build")
	if runtime.GOOS == "windows" {
		mainFolder = ".\\" + mainFolder
	} else {
		mainFolder = "./" + mainFolder
	}

	exp := []string{"go", "build", "-o", bn, mainFolder}

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

	args := []string{"builder", "a", "b", "c"}

	err := bc.Main(context.Background(), ".", args)
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

			plugs := plugins.Plugins{
				buildtest.BeforeBuilder(mainFileBeforeBuilder),
				buildtest.BeforeBuilder(fn),
				buildtest.AfterBuilder(cleanupAfterBuilder),
			}

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

			plugs := plugins.Plugins{
				buildtest.BeforeBuilder(mainFileBeforeBuilder),
				buildtest.AfterBuilder(cleanupAfterBuilder),
				buildtest.AfterBuilder(fn),
			}

			bc := &Cmd{
				pluginsFn: func() []plugins.Plugin {
					return plugs
				},
			}

			err := bc.Main(context.Background(), tt.root, tt.exp)
			if err != nil {
				r.Contains(err.Error(), tt.err.Error())
			}
			r.Equal(tt.exp, act)
		})
	}

}
