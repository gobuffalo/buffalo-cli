package generatecmd

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
)

type Command struct {
	Info      here.Info
	help      bool
	pluginsFn plugins.PluginFeeder
}

func (b *Command) WithHereInfo(i here.Info) {
	b.Info = i
}

func (b *Command) HereInfo() (here.Info, error) {
	if !b.Info.IsZero() {
		return b.Info, nil
	}
	return here.Current()
}

var _ plugins.PluginNeeder = &Command{}

func (b *Command) WithPlugins(f plugins.PluginFeeder) {
	b.pluginsFn = f
}

var _ plugprint.Aliases = &Command{}

func (*Command) Aliases() []string {
	return []string{"g"}
}

var _ plugins.Plugin = &Command{}

func (b Command) Name() string {
	return "generate"
}

func (b Command) String() string {
	return b.Name()
}

var _ plugprint.Describer = &Command{}

func (Command) Description() string {
	return "Generate application components"
}

var _ plugprint.SubCommander = &Command{}

func (bc *Command) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range bc.ScopedPlugins() {
		if _, ok := p.(Generator); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

var _ plugins.PluginScoper = &Command{}

func (bc *Command) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if bc.pluginsFn != nil {
		plugs = bc.pluginsFn()
	}

	var builders []plugins.Plugin
	for _, p := range plugs {
		switch p.(type) {
		case Generator:
			builders = append(builders, p)
		case BeforeGenerator:
			builders = append(builders, p)
		case AfterGenerator:
			builders = append(builders, p)
		}
	}
	return builders
}
