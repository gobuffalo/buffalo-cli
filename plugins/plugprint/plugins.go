package plugprint

import "github.com/gobuffalo/buffalo-cli/plugins"

type Plugins interface {
	WithPlugins() []plugins.Plugin
}
