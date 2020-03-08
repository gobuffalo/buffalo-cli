package newapp

import (
	"bytes"
	"context"
	"fmt"
	"os"
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
	fmt.Println(">>>TODO cli/cmds/newapp/cmd_test.go:43: info ", info)

	dir := filepath.Join(info.Root, "tmp")
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)
	fmt.Println(">>>TODO cli/cmds/newapp/cmd_test.go:48: dir ", dir)

	cmd := &Cmd{}

	pkg := "github.com/markbates/coke"
	name := "coke"
	err = cmd.Main(context.Background(), dir, []string{pkg})
	r.NoError(err)

	dir = filepath.Join(dir, name)
	mp := filepath.Join(dir, "go.mod")
	_, err = os.Stat(mp)
	r.NoError(err)

	info, err = here.Dir(dir)
	r.NoError(err)
	r.Equal(pkg, info.Module.Path)

	f, err := os.Open(filepath.Join(dir, "cmd", "newapp", "main.go"))
	r.NoError(err)
	r.NotNil(f)
}
