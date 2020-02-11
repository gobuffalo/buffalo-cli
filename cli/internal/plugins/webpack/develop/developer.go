package develop

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/develop"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack/internal/scripts"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"
)

var _ plugcmd.Namer = &Developer{}
var _ develop.Developer = &Developer{}
var _ plugins.Plugin = &Developer{}
var _ plugins.Needer = &Developer{}
var _ plugins.Scoper = &Developer{}

type Developer struct {
	pluginsFn plugins.Feeder
}

func (d *Developer) WithPlugins(fn plugins.Feeder) {
	d.pluginsFn = fn
}

func (d *Developer) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if d.pluginsFn == nil {
		return plugs
	}

	for _, p := range d.pluginsFn() {
		switch p.(type) {
		case Tooler:
			plugs = append(plugs, p)
		case Scripter:
			plugs = append(plugs, p)
		}
	}

	return plugs
}

func (d *Developer) PluginName() string {
	return "webpack/develop"
}

func (d *Developer) CmdName() string {
	return "webpack"
}

func (d *Developer) Develop(ctx context.Context, root string, args []string) error {
	tool, err := scripts.Tool(d, ctx, root)
	if err != nil {
		return err
	}

	// make sure that the node_modules folder is properly "installed"
	if _, err := os.Stat(filepath.Join(root, "node_modules")); err != nil {
		cmd := d.cmd(ctx, tool, "install")
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	cmd := d.cmd(ctx, scripts.WebpackBin(root), "--watch")

	if _, err := scripts.ScriptFor(d, ctx, root, "dev"); err == nil {
		cmd = d.cmd(ctx, tool, "run", "dev")
	}

	return cmd.Run()
}

func (d *Developer) cmd(ctx context.Context, tool string, args ...string) *exec.Cmd {
	plugs := d.ScopedPlugins()

	c := exec.CommandContext(ctx, tool, args...)
	c.Stdin = plugio.Stdin(plugs...)
	c.Stdout = plugio.Stdout(plugs...)
	c.Stderr = plugio.Stderr(plugs...)
	return c
}
