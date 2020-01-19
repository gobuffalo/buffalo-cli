package generate

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Cmd{},
	}
}

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

var _ plugins.PluginNeeder = &Cmd{}

func (b *Cmd) WithPlugins(f plugins.PluginFeeder) {
	b.pluginsFn = f
}

var _ plugprint.Aliases = &Cmd{}

func (*Cmd) Aliases() []string {
	return []string{"g"}
}

var _ plugins.Plugin = &Cmd{}

func (b Cmd) Name() string {
	return "generate"
}

func (b Cmd) String() string {
	return b.Name()
}

var _ plugprint.Describer = &Cmd{}

func (Cmd) Description() string {
	return "Generate application components"
}

var _ plugprint.SubCommander = &Cmd{}

func (bc *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range bc.ScopedPlugins() {
		if _, ok := p.(Generator); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

var _ plugins.PluginScoper = &Cmd{}

func (bc *Cmd) ScopedPlugins() []plugins.Plugin {
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
