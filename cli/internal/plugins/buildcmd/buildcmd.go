package buildcmd

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
)

var _ plugins.Plugin = &BuildCmd{}
var _ plugprint.Aliases = &BuildCmd{}
var _ plugprint.SubCommander = &BuildCmd{}
var _ plugprint.Describer = &BuildCmd{}
var _ plugprint.FlagPrinter = &BuildCmd{}
var _ plugprint.WithPlugins = &BuildCmd{}

type BuildCmd struct {
	plugins.IO
	Plugins func() plugins.Plugins
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
	dryRun                 bool
	help                   bool
	skipAssets             bool
	skipTemplateValidation bool
	verbose                bool
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

func (bc *BuildCmd) SubCommands() plugins.Plugins {
	var plugs plugins.Plugins
	for _, p := range bc.WithPlugins() {
		if _, ok := p.(Builder); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (bc *BuildCmd) WithPlugins() plugins.Plugins {
	var plugs plugins.Plugins
	if bc.Plugins != nil {
		plugs = bc.Plugins()
	}

	var builders plugins.Plugins
	for _, p := range plugs {
		switch p.(type) {
		case Builder:
			builders = append(builders, p)
		case BeforeBuilder:
			builders = append(builders, p)
		case AfterBuilder:
			builders = append(builders, p)
		case BuildVersioner:
			builders = append(builders, p)
		case TemplatesValidator:
			builders = append(builders, p)
		case Packager:
			builders = append(builders, p)
		}
	}
	return builders
}
