package plugprint

import "github.com/gobuffalo/buffalo-cli/plugins"

type Plugins interface {
	Plugins() []plugins.Plugin
}
