package testcmd

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
)

type TestCmd struct {
	Info      here.Info
	pluginsFn plugins.PluginFeeder
}

var _ plugins.Plugin = &TestCmd{}

func (TestCmd) Name() string {
	return "test"
}

var _ plugprint.Describer = &TestCmd{}

func (TestCmd) Description() string {
	return "Run the tests for the Buffalo app."
}

func (b *TestCmd) WithHereInfo(i here.Info) {
	b.Info = i
}

func (b *TestCmd) HereInfo() (here.Info, error) {
	if !b.Info.IsZero() {
		return b.Info, nil
	}
	return here.Current()
}

var _ plugins.PluginNeeder = &TestCmd{}

func (b *TestCmd) WithPlugins(f plugins.PluginFeeder) {
	b.pluginsFn = f
}

var _ plugins.PluginScoper = &TestCmd{}

func (bc *TestCmd) ScopedPlugins() []plugins.Plugin {
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

var _ plugprint.SubCommander = &TestCmd{}

func (bc *TestCmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range bc.ScopedPlugins() {
		if _, ok := p.(Tester); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}
