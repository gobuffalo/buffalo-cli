package newapp

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

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
	orig := []string{"-p", "json", "-p", "github.com/other/preset", pkg}

	err = cmd.Main(context.Background(), dir, orig)
	r.NoError(err)

	exp := []string{"go", "run", "./cmd/newapp"}
	exp = append(exp, orig[:len(orig)-1]...)
	r.Equal(exp, act)

	dir = filepath.Join(dir, name)
	mp := filepath.Join(dir, "go.mod")
	_, err = os.Stat(mp)
	r.NoError(err)

	b, err := ioutil.ReadFile(filepath.Join(dir, "cmd", "newapp", "main.go"))
	r.NoError(err)

	ba := string(b)
	ba = strings.TrimSpace(ba)
	be := strings.TrimSpace(newappExp)
	r.Equal(be, ba)
}

const newappExp = `
package main

import (
	"context"
	"log"
	"os"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp"
	"github.com/gobuffalo/plugins"

	json "github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp/presets/jsonapp"
	preset "github.com/other/preset"
)

func main() {
	ctx := context.Background()
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var plugs []plugins.Plugin

	plugs = append(plugs, json.Plugins()...)
	plugs = append(plugs, preset.Plugins()...)

	args := []string{"-p", "json", "-p", "github.com/other/preset"}
	if err := newapp.Execute(plugs, ctx, pwd, "github.com/markbates/coke", args); err != nil {
		log.Fatal(err)
	}
}
`
