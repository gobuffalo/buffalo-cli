package setup

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/spf13/pflag"
)

var _ plugcmd.Commander = &Cmd{}
var _ plugcmd.SubCommander = &Cmd{}
var _ plugins.Needer = &Cmd{}
var _ plugins.Plugin = &Cmd{}
var _ plugins.Scoper = &Cmd{}

type Cmd struct {
	pluginsFn plugins.Feeder
	flags     *pflag.FlagSet
	help      bool
}

func (Cmd) PluginName() string {
	return "setup"
}

func (cmd *Cmd) WithPlugins(f plugins.Feeder) {
	cmd.pluginsFn = f
}

func (cmd *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if cmd.pluginsFn == nil {
		return plugs
	}

	for _, p := range cmd.pluginsFn() {
		switch p.(type) {
		case Setuper:
			plugs = append(plugs, p)
		case BeforeSetuper:
			plugs = append(plugs, p)
		case AfterSetuper:
			plugs = append(plugs, p)
		case Stdouter:
			plugs = append(plugs, p)
		case Flagger:
			plugs = append(plugs, p)
		case Pflagger:
			plugs = append(plugs, p)
		}
	}

	return plugs
}

func (cmd *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin

	for _, p := range cmd.ScopedPlugins() {
		switch p.(type) {
		case Setuper:
			plugs = append(plugs, p)
		}
	}

	return plugs
}
