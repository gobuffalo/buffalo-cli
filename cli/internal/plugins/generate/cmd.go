package generate

import (
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &Cmd{}
var _ plugins.PluginNeeder = &Cmd{}
var _ plugins.PluginScoper = &Cmd{}
var _ plugprint.Aliases = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}
var _ plugprint.SubCommander = &Cmd{}

type Cmd struct {
	help      bool
	pluginsFn plugins.PluginFeeder
	flags     *pflag.FlagSet
}

func (cmd *Cmd) WithPlugins(f plugins.PluginFeeder) {
	cmd.pluginsFn = f
}

func (*Cmd) Aliases() []string {
	return []string{"g"}
}

func (b Cmd) Name() string {
	return "generate"
}

func (Cmd) Description() string {
	return "Generate application components"
}

func (cmd *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range cmd.ScopedPlugins() {
		if _, ok := p.(Generator); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (cmd *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if cmd.pluginsFn == nil {
		return plugs
	}

	for _, p := range cmd.pluginsFn() {
		switch p.(type) {
		case Generator:
			plugs = append(plugs, p)
		}
	}
	return plugs
}
