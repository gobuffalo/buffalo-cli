package resource

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Generator{},
	}
}

type Generator struct {
	info               here.Info
	pluginsFn          plugins.PluginFeeder
	help               bool
	skipActionTests    bool
	skipActions        bool
	skipMigrationTests bool
	skipMigrations     bool
	skipModelTests     bool
	skipModels         bool
	skipTemplateTests  bool
	skipTemplates      bool
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

	pm := map[string]bool{}

	var builders []plugins.Plugin

	for _, p := range plugs {
		_, ok := p.(BeforeGenerator)
		if !ok {
			continue
		}
		builders = append(builders, p)
	}

	for _, p := range plugs {
		switch p.(type) {
		case ResourceGenerator:
			if pm["ResourceGenerator"] {
				continue
			}
			pm["ResourceGenerator"] = true
			break
		case Actioner:
			if pm["Actioner"] {
				continue
			}
			pm["Actioner"] = true
		case ActionTester:
			if pm["ActionTester"] {
				continue
			}
			pm["ActionTester"] = true
		case Modeler:
			if pm["Modeler"] {
				continue
			}
			pm["Modeler"] = true
		case ModelTester:
			if pm["ModelTester"] {
				continue
			}
			pm["ModelTester"] = true
		case Templater:
			if pm["Templater"] {
				continue
			}
			pm["Templater"] = true
		case TemplateTester:
			if pm["TemplateTester"] {
				continue
			}
			pm["TemplateTester"] = true
		default:
			continue
		}

		builders = append(builders, p)
	}

	for _, p := range plugs {
		_, ok := p.(AfterGenerator)
		if !ok {
			continue
		}
		builders = append(builders, p)
	}

	return builders
}
