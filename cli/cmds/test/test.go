package test

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugprint"
)

var _ plugins.Plugin = &Cmd{}
var _ plugins.Needer = &Cmd{}
var _ plugins.Scoper = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.SubCommander = &Cmd{}

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Cmd{},
		&Setup{},
	}
}

type Cmd struct {
	pluginsFn plugins.Feeder
}

func (Cmd) PluginName() string {
	return "test"
}

func (Cmd) Description() string {
	return "Run the tests for the Buffalo app."
}

func (b *Cmd) WithPlugins(f plugins.Feeder) {
	b.pluginsFn = f
}

func (bc *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if bc.pluginsFn != nil {
		plugs = bc.pluginsFn()
	}

	var builders []plugins.Plugin
	for _, p := range plugs {
		switch p.(type) {
		case Tester:
			builders = append(builders, p)
		case BeforeTester:
			builders = append(builders, p)
		case AfterTester:
			builders = append(builders, p)
		case Runner:
			builders = append(builders, p)
		case Argumenter:
			builders = append(builders, p)
		case Stdouter:
			builders = append(builders, p)
		case Namer:
			builders = append(builders, p)
		case Aliaser:
			builders = append(builders, p)
		}
	}
	return builders
}

func (bc *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range bc.ScopedPlugins() {
		if _, ok := p.(Tester); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}
