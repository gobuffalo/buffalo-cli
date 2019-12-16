package plugprint

import "github.com/gobuffalo/buffalo-cli/plugins"

type WithPlugins interface {
	WithPlugins() []plugins.Plugin
}
