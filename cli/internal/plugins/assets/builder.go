package assets

import (
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
)

var _ buildcmd.BeforeBuilder = &Builder{}
var _ buildcmd.BuildPflagger = &Builder{}
var _ plugins.Plugin = &Builder{}
var _ plugprint.Describer = &Builder{}
var _ plugprint.FlagPrinter = &Builder{}

// Builder is responsible for building webpack
// and other assets
type Builder struct {
	Environment string
	// CleanAssets will remove the public/assets folder build compiling
	Clean bool
	// places ./public/assets into ./bin/assets.zip.
	Extract   bool
	ExtractTo string // ./bin

	AssetPaths []string

	Skip bool
	Tool string // default is npm
}

func (a Builder) Name() string {
	return "assets"
}

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
