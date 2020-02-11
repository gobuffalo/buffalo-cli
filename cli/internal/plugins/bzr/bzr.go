package bzr

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/build"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugprint"
)

var _ build.Versioner = &Versioner{}
var _ plugins.Plugin = Versioner{}
var _ plugins.Needer = &Versioner{}
var _ plugprint.Describer = Versioner{}

// Plugins ...
func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		Versioner{},
	}
}
