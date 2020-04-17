package newapp

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Help(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}
	stdout := plugio.NewOuter(bb)

	cmd := &Cmd{}
	cmd.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			stdout,
		}
	})

	err := cmd.Main(context.Background(), "", []string{"-h"})
	r.NoError(err)

	body := bb.String()
	r.Contains(body, `$ new`)
}

func Test_Cmd_Main(t *testing.T) {
	r := require.New(t)

	info, err := here.Current()
	r.NoError(err)

	dir := filepath.Join(info.Root, "tmp")
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	cmd := &Cmd{}
	var act []string
	fn := func(ctx context.Context, root string, cmd *exec.Cmd) error {
		act = cmd.Args
		return nil
	}
	cmd.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{
			cmdRunner(fn),
		}
	})

	pkg := "github.com/markbates/coke"
	name := "coke"
	err = cmd.Main(context.Background(), dir, []string{pkg})
	r.NoError(err)

	exp := []string{"go", "run", "./cmd/newapp"}
	r.Equal(exp, act)

	dir = filepath.Join(dir, name)
	mp := filepath.Join(dir, "go.mod")
	_, err = os.Stat(mp)
	r.NoError(err)

	f, err := os.Open(filepath.Join(dir, "cmd", "newapp", "main.go"))
	r.NoError(err)
	r.NotNil(f)
}
