package info

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &Cmd{}
var _ plugins.Needer = &Cmd{}
var _ plugins.Scoper = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}

type Cmd struct {
	flags     *pflag.FlagSet
	pluginsFn plugins.Feeder
	help      bool
}

func (cmd *Cmd) WithPlugins(f plugins.Feeder) {
	cmd.pluginsFn = f
}

func (cmd *Cmd) PluginName() string {
	return "info"
}

func (cmd *Cmd) Description() string {
	return "Print diagnostic information (useful for debugging)"
}

func (cmd *Cmd) ScopedPlugins() []plugins.Plugin {
	if cmd.pluginsFn == nil {
		return nil
	}

	var plugs []plugins.Plugin

	for _, p := range cmd.pluginsFn() {
		switch p.(type) {
		case Informer:
			plugs = append(plugs, p)
		}
	}
	return plugs
}
