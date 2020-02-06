package develop

import (
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/spf13/pflag"
)

type Cmd struct {
	pluginsFn plugins.PluginFeeder
	flags     *pflag.FlagSet
	help      bool
}

func (cmd *Cmd) WithPlugins(f plugins.PluginFeeder) {
	cmd.pluginsFn = f
}

func (cmd *Cmd) ScopedPlugins() []plugins.Plugin {
	if cmd.pluginsFn == nil {
		return []plugins.Plugin{}
	}

	var plugs []plugins.Plugin
	for _, p := range cmd.pluginsFn() {
		switch p.(type) {
		case Developer:
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (cmd *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range cmd.pluginsFn() {
		switch p.(type) {
		case Developer:
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (cmd *Cmd) PluginName() string {
	return "develop/cmd"
}

func (cmd *Cmd) CmdName() string {
	return "develop"
}

func (cmd *Cmd) Aliases() []string {
	return []string{"dev"}
}

func (cmd *Cmd) Description() string {
	return "run go apps in 'development' mode"
}
