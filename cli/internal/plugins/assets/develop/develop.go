package develop

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/assets/scripts"
	"github.com/gobuffalo/buffalo-cli/plugins"
)

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

func (d *Developer) Name() string {
	return "assets/develop"
}

func (d *Developer) CmdName() string {
	return "assets"
}

func (d *Developer) Develop(ctx context.Context, root string, args []string) error {
	tool, err := d.tool(ctx, root)
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

	if _, err := d.scriptFor(ctx, root, "dev"); err == nil {
		cmd = plugins.Cmd(ctx, tool, "run", "dev")
	}

	return cmd.Run()
}

func (d *Developer) scriptFor(ctx context.Context, root string, name string) (string, error) {
	for _, p := range d.ScopedPlugins() {
		if tp, ok := p.(Scripter); ok {
			return tp.AssetScript(ctx, root, name)
		}
	}
	return scripts.ScriptFor(root, name)
}

func (d *Developer) tool(ctx context.Context, root string) (string, error) {
	for _, p := range d.ScopedPlugins() {
		if tp, ok := p.(Tooler); ok {
			return tp.AssetTool(ctx, root)
		}
	}
	return scripts.Tool(root)
}
