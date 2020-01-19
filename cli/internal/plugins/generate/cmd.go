package generate

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
)

var _ plugins.Plugin = &Cmd{}
var _ plugins.PluginNeeder = &Cmd{}
var _ plugins.PluginScoper = &Cmd{}
var _ plugprint.Aliases = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}
var _ plugprint.SubCommander = &Cmd{}

type Cmd struct {
	Info      here.Info
	help      bool
	pluginsFn plugins.PluginFeeder
}

func (b *Cmd) WithHereInfo(i here.Info) {
	b.Info = i
}

func (b *Cmd) HereInfo() (here.Info, error) {
	if !b.Info.IsZero() {
		return b.Info, nil
	}
	return here.Current()
}

func (b *Cmd) WithPlugins(f plugins.PluginFeeder) {
	b.pluginsFn = f
}

func (*Cmd) Aliases() []string {
	return []string{"g"}
}

func (b Cmd) Name() string {
	return "generate"
}

func (b Cmd) String() string {
	return b.Name()
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
