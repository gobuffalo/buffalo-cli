package develop

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack/internal/scripts"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/develop"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
)

var _ develop.Developer = &Developer{}
var _ plugins.NamedCommand = &Developer{}
var _ plugins.Plugin = &Developer{}
var _ plugins.PluginNeeder = &Developer{}
var _ plugins.PluginScoper = &Developer{}

type Developer struct {
	pluginsFn plugins.PluginFeeder
}

func (d *Developer) WithPlugins(fn plugins.PluginFeeder) {
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
		cmd := plugins.Cmd(ctx, tool, "install")
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	cmd := plugins.Cmd(ctx, scripts.WebpackBin(root), "--watch")

	if _, err := scripts.ScriptFor(d, ctx, root, "dev"); err == nil {
		cmd = plugins.Cmd(ctx, tool, "run", "dev")
	}

	return cmd.Run()
}
