package build

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ build.BeforeBuilder = &Builder{}
var _ build.Builder = &Builder{}
var _ build.Pflagger = &Builder{}
var _ plugcmd.Namer = &Builder{}
var _ plugins.Needer = &Builder{}
var _ plugins.Plugin = &Builder{}
var _ plugins.Scoper = &Builder{}
var _ plugprint.Describer = &Builder{}
var _ plugprint.FlagPrinter = &Builder{}

// Builder is responsible for building webpack
type Builder struct {
	environment string
	clean       bool   // CleanAssets will remove the public/webpack folder build compiling
	extract     bool   // places ./public/webpack into ./bin/webpack.zip.
	extractTo   string // ./bin
	assetPaths  []string
	skip        bool
	tool        string // default is npm
	pluginsFn   plugins.Feeder
	flags       *pflag.FlagSet
}

func (bc *Builder) WithPlugins(f plugins.Feeder) {
	bc.pluginsFn = f
}

func (bc *Builder) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if bc.pluginsFn != nil {
		plugs = bc.pluginsFn()
	}

	var builders []plugins.Plugin
	for _, p := range plugs {
		switch p.(type) {
		case AssetBuilder:
			builders = append(builders, p)
		case Tooler:
			builders = append(builders, p)
		case Scripter:
			builders = append(builders, p)
		case Stdouter:
			builders = append(builders, p)
		}
	}
	return builders
}

func (a Builder) PluginName() string {
	return "webpack/builder"
}

func (a Builder) CmdName() string {
	return "webpack"
}

func (a Builder) Description() string {
	return "Builds webpack applications"
}

func (a Builder) String() string {
	return a.PluginName()
}
