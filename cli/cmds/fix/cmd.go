package fix

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugcmd.SubCommander = &Cmd{}
var _ plugins.Needer = &Cmd{}
var _ plugins.Plugin = &Cmd{}
var _ plugins.Scoper = &Cmd{}
var _ plugprint.Describer = &Cmd{}

type Cmd struct {
	flags     *pflag.FlagSet
	help      bool
	pluginsFn plugins.Feeder
}

func (fc *Cmd) WithPlugins(f plugins.Feeder) {
	fc.pluginsFn = f
}

func (fc *Cmd) PluginName() string {
	return "fix"
}

func (fc *Cmd) Description() string {
	return "Attempt to fix a Buffalo application's API to match version in go.mod"
}

func (fc *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range fc.ScopedPlugins() {
		if c, ok := p.(Fixer); ok {
			plugs = append(plugs, c)
		}
	}
	return plugs
}

func (fc *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if fc.pluginsFn == nil {
		return plugs
	}

	for _, p := range fc.pluginsFn() {
		switch p.(type) {
		case Fixer:
			plugs = append(plugs, p)
		case Flagger:
			plugs = append(plugs, p)
		case Pflagger:
			plugs = append(plugs, p)
		case Stdouter:
			plugs = append(plugs, p)
		}
	}
	return plugs
}
