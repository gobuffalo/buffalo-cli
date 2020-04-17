package newapp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"text/template"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp/presets"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/cligen"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &Cmd{}
var _ plugcmd.Commander = &Cmd{}
var _ plugcmd.Namer = &Cmd{}
var _ plugins.Needer = &Cmd{}
var _ plugins.Scoper = &Cmd{}

type Cmd struct {
	pluginsFn plugins.Feeder
	flags     *pflag.FlagSet
	help      bool
	force     bool
	preset    string
	usePlugs  map[string]string
}

func (Cmd) PluginName() string {
	return "newapp/cmd"
}

func (Cmd) CmdName() string {
	return "new"
}

func (cmd *Cmd) WithPlugins(f plugins.Feeder) {
	cmd.pluginsFn = f
}

func (cmd *Cmd) ScopedPlugins() []plugins.Plugin {
	if cmd.pluginsFn == nil {
		return nil
	}

	var plugs []plugins.Plugin

	for _, p := range cmd.pluginsFn() {
		switch p.(type) {
		case Stdouter:
			plugs = append(plugs, p)
		case Stdiner:
			plugs = append(plugs, p)
		case Stderrer:
			plugs = append(plugs, p)
		case NewCommandRunner:
			plugs = append(plugs, p)
		}
	}

	return plugs
}

func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return plugins.Wrap(cmd, err)
	}

	plugs := cmd.ScopedPlugins()

	if cmd.help {
		return plugprint.Print(plugio.Stdout(plugs...), cmd)
	}

	args = flags.Args()

	if len(args) == 0 {
		return plugins.Wrap(cmd, fmt.Errorf("missing application name"))
	}

	modName := args[0]
	dirName := path.Base(modName)
	args = args[1:]

	root = filepath.Join(root, dirName)
	if cmd.force {
		os.RemoveAll(root)
	}

	if err := os.MkdirAll(root, 0755); err != nil {
		return plugins.Wrap(cmd, err)
	}

	os.Chdir(root)

	if err := cmd.modInit(ctx, root, modName); err != nil {
		return plugins.Wrap(cmd, err)
	}

	if cmd.usePlugs == nil {
		cmd.usePlugs = map[string]string{}
	}

	if len(cmd.preset) > 0 {
		pres := presets.Presets()
		for _, p := range pres {
			if path.Base(p) != cmd.preset {
				continue
			}
			cmd.usePlugs[cmd.preset] = p
			break
		}
	}

	tmpl, err := template.New("").Parse(cliMain)
	if err != nil {
		return plugins.Wrap(cmd, err)
	}

	cd := filepath.Join(root, "cmd", "newapp")
	if err := os.MkdirAll(cd, 0755); err != nil {
		return plugins.Wrap(cmd, err)
	}

	w, err := os.Create(filepath.Join(cd, "main.go"))
	if err != nil {
		return err
	}
	defer w.Close()

	err = tmpl.Execute(w, map[string]interface{}{
		"Plugs": cmd.usePlugs,
	})
	if err != nil {
		return plugins.Wrap(cmd, err)
	}

	g := &cligen.Generator{
		Plugins: cmd.usePlugs,
	}
	if err := g.Generate(ctx, root, args); err != nil {
		return plugins.Wrap(cmd, err)
	}

	os.Chdir(root)

	c := exec.CommandContext(ctx, "go", "run", "./cmd/newapp")
	c.Stdout = plugio.Stdout(plugs...)
	c.Stderr = plugio.Stderr(plugs...)
	c.Stdin = plugio.Stdin(plugs...)

	for _, p := range plugs {
		if vr, ok := p.(NewCommandRunner); ok {
			if err := vr.RunNewCommand(ctx, root, c); err != nil {
				return plugins.Wrap(cmd, err)
			}
			return nil
		}
	}

	bb := &bytes.Buffer{}
	c.Stderr = io.MultiWriter(bb, c.Stderr)
	if err := c.Run(); err != nil {
		return plugins.Wrap(cmd, fmt.Errorf("%w: %s", err, bb.String()))
	}
	return nil
}

const cliMain = `
package main

import (
	"context"
	"log"
	"os"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp"
	"github.com/gobuffalo/plugins"
	{{range $k,$v := .Plugs -}}
	"{{$v}}"
	{{- end}}
)

func main() {
	ctx := context.Background()
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var plugs []plugins.Plugin
	{{range $k,$v := .Plugs -}}
	plugs = append(plugs, {{$k}}.Plugins()...)
	{{- end}}

	if err := newapp.Execute(plugs, ctx, pwd, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
`
