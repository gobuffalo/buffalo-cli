package build

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"github.com/gobuffalo/here"
	"github.com/spf13/pflag"
)

var _ build.BeforeBuilder = &Builder{}
var _ build.Builder = &Builder{}
var _ build.Pflagger = &Builder{}
var _ plugins.NamedCommand = &Builder{}
var _ plugins.Plugin = &Builder{}
var _ plugins.PluginNeeder = &Builder{}
var _ plugins.PluginScoper = &Builder{}
var _ plugprint.Describer = &Builder{}
var _ plugprint.FlagPrinter = &Builder{}

// Builder is responsible for building webpack
// and other webpack
type Builder struct {
	Environment string
	// CleanAssets will remove the public/webpack folder build compiling
	Clean bool
	// places ./public/webpack into ./bin/webpack.zip.
	Extract    bool
	ExtractTo  string // ./bin
	AssetPaths []string
	Skip       bool
	Tool       string // default is npm

	info      here.Info
	pluginsFn plugins.PluginFeeder
	flags     *pflag.FlagSet
}

func (bc *Builder) WithPlugins(f plugins.PluginFeeder) {
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
		}
	}
	return builders
}

func (a *Builder) WithHereInfo(i here.Info) {
	a.info = i
}

func (a *Builder) HereInfo() (here.Info, error) {
	if !a.info.IsZero() {
		return a.info, nil
	}
	return here.Current()
}

func (a Builder) Name() string {
	return "webpack/builder"
}

func (a Builder) CmdName() string {
	return "webpack"
}

func (a Builder) Description() string {
	return "Builds webpack applications"
}

func (a Builder) String() string {
	return a.Name()
}
