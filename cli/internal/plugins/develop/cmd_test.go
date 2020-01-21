package develop

import (
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
)

var _ plugins.Aliases = &Cmd{}
var _ plugins.NamedCommand = &Cmd{}
var _ plugins.Plugin = &Cmd{}
var _ plugins.PluginNeeder = &Cmd{}
var _ plugins.PluginScoper = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}
var _ plugprint.SubCommander = &Cmd{}
