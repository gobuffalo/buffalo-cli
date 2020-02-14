package grifts

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/grifts/cmd"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/grifts/generator"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/grifts/importer"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/grifts/setup"
	"github.com/gobuffalo/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&cmd.Cmd{},
		&generator.Generator{},
		&importer.Importer{},
		&setup.Setup{},
	}
}
