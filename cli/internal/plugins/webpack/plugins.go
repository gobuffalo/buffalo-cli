package webpack

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack/build"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack/develop"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
)

// Plugins returns all of the plugins available in this package.
// All plugins use zero values.
func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&build.Builder{},
		&develop.Developer{},
	}
}
