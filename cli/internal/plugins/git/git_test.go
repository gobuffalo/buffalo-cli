package git

import (
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
)

var _ build.Versioner = &Versioner{}
var _ plugins.Plugin = Versioner{}
var _ plugins.PluginNeeder = &Versioner{}
var _ plugprint.Describer = Versioner{}
