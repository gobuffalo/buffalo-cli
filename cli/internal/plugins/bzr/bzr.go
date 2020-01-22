package bzr

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
)

var _ buildcmd.Versioner = &Versioner{}
var _ plugins.Plugin = Versioner{}
var _ plugins.PluginNeeder = &Versioner{}
var _ plugprint.Describer = Versioner{}

// Plugins ...
func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		Versioner{},
	}
}
