package test

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
	pluginsFn plugins.PluginFeeder
}

var _ plugins.Plugin = &Cmd{}

func (Cmd) Name() string {
	return "test"
}

var _ plugprint.Describer = &Cmd{}

func (Cmd) Description() string {
	return "Run the tests for the Buffalo app."
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

var _ plugins.PluginScoper = &Cmd{}

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
		}
	}
	return builders
}

var _ plugprint.SubCommander = &Cmd{}

func (bc *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range bc.ScopedPlugins() {
		if _, ok := p.(Tester); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}
