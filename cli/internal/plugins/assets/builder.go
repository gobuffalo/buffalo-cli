package assets

import (
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/plugins/plugprint"
	"github.com/gobuffalo/here"
)

// Builder is responsible for building webpack
// and other assets
type Builder struct {
	Info here.Info

	Environment string
	// CleanAssets will remove the public/assets folder build compiling
	Clean bool
	// places ./public/assets into ./bin/assets.zip.
	Extract   bool
	ExtractTo string // ./bin

	AssetPaths []string

	Skip bool
	Tool string // default is npm

	pluginsFn plugins.PluginFeeder
}

var _ plugins.PluginNeeder = &Builder{}

func (bc *Builder) WithPlugins(f plugins.PluginFeeder) {
	bc.pluginsFn = f
}

var _ plugins.PluginScoper = &Builder{}

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
		}
	}
	return builders
}

func (a *Builder) WithHereInfo(i here.Info) {
	a.Info = i
}

func (a *Builder) HereInfo() (here.Info, error) {
	if !a.Info.IsZero() {
		return a.Info, nil
	}
	return here.Current()
}

var _ plugins.Plugin = &Builder{}

func (a Builder) Name() string {
	return "assets"
}

var _ plugprint.Describer = &Builder{}

func (a Builder) Description() string {
	return "Manages webpack assets during the buffalo build process."
}

func (a Builder) String() string {
	return a.Name()
}

func (b Builder) webpackBin() string {
	s := filepath.Join("node_modules", ".bin", "webpack")
	if runtime.GOOS == "windows" {
		s += ".cmd"
	}
	return s
}
