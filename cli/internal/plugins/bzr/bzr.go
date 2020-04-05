package bzr

import (
	"github.com/gobuffalo/plugins"
)

// Plugins ...
func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		Versioner{},
	}
}
