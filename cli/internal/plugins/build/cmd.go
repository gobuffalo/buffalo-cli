package build

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
)

type Cmd struct {
	Info here.Info

	// Mod is the -mod flag
	Mod string
	// Static sets the following flags for the final `go build` command:
	// -linkmode external
	// -extldflags "-static"
	Static bool
	// Environment the binary is meant for. defaults to "development"
	Environment string
	// LDFlags to be passed to the final `go build` command
	LDFlags string
	// BuildFlags to be passed to the final `go build` command
	BuildFlags             []string
	Tags                   string
	Bin                    string
	help                   bool
	skipTemplateValidation bool
	verbose                bool

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
	return []string{"b", "install"}
}

var _ plugins.Plugin = &Cmd{}

func (b Cmd) Name() string {
	return "build"
}

func (b Cmd) String() string {
	return b.Name()
}

var _ plugprint.Describer = &Cmd{}

func (Cmd) Description() string {
	return "Build the application binary, including bundling of assets (packr & webpack)"
}

var _ plugprint.SubCommander = &Cmd{}

func (bc *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range bc.ScopedPlugins() {
		if _, ok := p.(Builder); ok {
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
		case Builder:
			builders = append(builders, p)
		case BeforeBuilder:
			builders = append(builders, p)
		case AfterBuilder:
			builders = append(builders, p)
		case Versioner:
			builders = append(builders, p)
		case TemplatesValidator:
			builders = append(builders, p)
		case Packager:
			builders = append(builders, p)
		case PackFiler:
			builders = append(builders, p)
		case Flagger:
			builders = append(builders, p)
		case Pflagger:
			builders = append(builders, p)
		case Importer:
			builders = append(builders, p)
		case Runner:
			builders = append(builders, p)
		}
	}
	return builders
}
