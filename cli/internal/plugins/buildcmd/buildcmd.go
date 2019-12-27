package buildcmd

import (
	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/plugins/plugprint"
	"github.com/gobuffalo/here"
)

var _ plugins.Plugin = &BuildCmd{}
var _ plugins.PluginNeeder = &BuildCmd{}
var _ plugins.PluginScoper = &BuildCmd{}
var _ plugprint.Aliases = &BuildCmd{}
var _ plugprint.Describer = &BuildCmd{}
var _ plugprint.FlagPrinter = &BuildCmd{}
var _ plugprint.SubCommander = &BuildCmd{}

type BuildCmd struct {
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

func (b *BuildCmd) WithHereInfo(i here.Info) {
	b.Info = i
}

func (b *BuildCmd) HereInfo() (here.Info, error) {
	if !b.Info.IsZero() {
		return b.Info, nil
	}
	return here.Current()
}

func (b *BuildCmd) WithPlugins(f plugins.PluginFeeder) {
	b.pluginsFn = f
}

func (*BuildCmd) Aliases() []string {
	return []string{"b", "install"}
}

func (b BuildCmd) Name() string {
	return "build"
}

func (b BuildCmd) String() string {
	return b.Name()
}

func (BuildCmd) Description() string {
	return "Build the application binary, including bundling of assets (packr & webpack)"
}

func (bc *BuildCmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range bc.ScopedPlugins() {
		if _, ok := p.(Builder); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (bc *BuildCmd) ScopedPlugins() []plugins.Plugin {
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
