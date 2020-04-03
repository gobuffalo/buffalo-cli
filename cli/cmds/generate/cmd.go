package generate

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugcmd.Aliaser = &Cmd{}
var _ plugins.Needer = &Cmd{}
var _ plugins.Plugin = &Cmd{}
var _ plugins.Scoper = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}
var _ plugprint.SubCommander = &Cmd{}

type Cmd struct {
	help      bool
	pluginsFn plugins.Feeder
	flags     *pflag.FlagSet
}

func (cmd *Cmd) WithPlugins(f plugins.Feeder) {
	cmd.pluginsFn = f
}

func (*Cmd) CmdAliases() []string {
	return []string{"g"}
}

func (b Cmd) PluginName() string {
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
		case Stdouter:
			plugs = append(plugs, p)
		}
	}
	return plugs
}
