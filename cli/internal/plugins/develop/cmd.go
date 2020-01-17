package develop

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here"
	"github.com/spf13/pflag"
)

type Cmd struct {
	pluginsFn plugins.PluginFeeder
	info      here.Info
	flags     *pflag.FlagSet
	help      bool
}

func (cmd *Cmd) WithHereInfo(i here.Info) {
	cmd.info = i
}

func (cmd *Cmd) HereInfo() (here.Info, error) {
	if cmd.info.IsZero() {
		return here.Current()
	}
	return cmd.info, nil
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

func (cmd *Cmd) Name() string {
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
