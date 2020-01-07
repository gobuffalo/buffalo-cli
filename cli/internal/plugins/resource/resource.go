package resource

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here"
)

type Generator struct {
	info      here.Info
	pluginsFn plugins.PluginFeeder
	help      bool
}

var _ plugins.PluginNeeder = &Generator{}

func (g *Generator) WithPlugins(f plugins.PluginFeeder) {
	g.pluginsFn = f
}

func (g *Generator) WithHereInfo(i here.Info) {
	g.info = i
}

func (g *Generator) HereInfo() (here.Info, error) {
	if g.info.IsZero() {
		return here.Current()
	}
	return g.info, nil
}

var _ plugins.Plugin = Generator{}

func (g Generator) Name() string {
	return "resource"
}

var _ plugins.Aliases = Generator{}

func (g Generator) Aliases() []string {
	return []string{"r"}
}

var _ plugins.PluginScoper = &Generator{}

func (g *Generator) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if g.pluginsFn != nil {
		plugs = g.pluginsFn()
	}

	var builders []plugins.Plugin
	for _, p := range plugs {
		switch p.(type) {
		case BeforeGenerator:
			builders = append(builders, p)
		case ResourceGenerator:
			builders = append(builders, p)
		case AfterGenerator:
			builders = append(builders, p)
		case Actioner:
			builders = append(builders, p)
		case ActionTester:
			builders = append(builders, p)
		case Modeler:
			builders = append(builders, p)
		case ModelTester:
			builders = append(builders, p)
		case Migrationer:
			builders = append(builders, p)
		case MigrationTester:
			builders = append(builders, p)
		case Templater:
			builders = append(builders, p)
		case TemplateTester:
			builders = append(builders, p)
		}
	}
	return builders
}
