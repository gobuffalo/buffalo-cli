package build

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &Cmd{}
var _ plugins.PluginNeeder = &Cmd{}
var _ plugins.PluginScoper = &Cmd{}
var _ plugprint.Aliases = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}
var _ plugprint.SubCommander = &Cmd{}

type Cmd struct {
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
	Verbose                bool
	SkipTemplateValidation bool

	help      bool
	info      here.Info
	pluginsFn plugins.PluginFeeder
	flags     *pflag.FlagSet
}

func (cmd *Cmd) WithHereInfo(i here.Info) {
	cmd.info = i
}

func (cmd *Cmd) HereInfo() (here.Info, error) {
	if !cmd.info.IsZero() {
		return cmd.info, nil
	}
	return here.Current()
}

func (cmd *Cmd) WithPlugins(f plugins.PluginFeeder) {
	cmd.pluginsFn = f
}

func (*Cmd) Aliases() []string {
	return []string{"b", "install"}
}

func (b Cmd) Name() string {
	return "build"
}

func (b Cmd) String() string {
	return b.Name()
}

func (Cmd) Description() string {
	return "Build the application binary, including bundling of assets (packr & webpack)"
}

func (bc *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range bc.ScopedPlugins() {
		if _, ok := p.(Builder); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

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
