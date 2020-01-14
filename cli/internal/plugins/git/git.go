package git

import (
	"github.com/gobuffalo/buffalo-cli/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Versioner{},
	}
}
