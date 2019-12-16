package plugprint

import "github.com/gobuffalo/buffalo-cli/plugins"

type SubCommander interface {
	SubCommands() []plugins.Plugin
}
