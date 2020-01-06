package generatecmd

import (
	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/plugins/plugprint"
	"github.com/gobuffalo/here"
)

type GenerateCmd struct {
	Info      here.Info
	help      bool
	pluginsFn plugins.PluginFeeder
}

func (b *GenerateCmd) WithHereInfo(i here.Info) {
	b.Info = i
}

func (b *GenerateCmd) HereInfo() (here.Info, error) {
	if !b.Info.IsZero() {
		return b.Info, nil
	}
	return here.Current()
}

var _ plugins.PluginNeeder = &GenerateCmd{}

func (b *GenerateCmd) WithPlugins(f plugins.PluginFeeder) {
	b.pluginsFn = f
}

var _ plugprint.Aliases = &GenerateCmd{}

func (*GenerateCmd) Aliases() []string {
	return []string{"g"}
}

var _ plugins.Plugin = &GenerateCmd{}

func (b GenerateCmd) Name() string {
	return "generate"
}

func (b GenerateCmd) String() string {
	return b.Name()
}

var _ plugprint.Describer = &GenerateCmd{}

func (GenerateCmd) Description() string {
	return "Generate application components"
}

var _ plugprint.SubCommander = &GenerateCmd{}

func (bc *GenerateCmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range bc.ScopedPlugins() {
		if _, ok := p.(Generator); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

var _ plugins.PluginScoper = &GenerateCmd{}

func (bc *GenerateCmd) ScopedPlugins() []plugins.Plugin {
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
