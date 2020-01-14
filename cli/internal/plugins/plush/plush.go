package plush

import (
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/plush/internal/generators/resource"
	"github.com/gobuffalo/buffalo-cli/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Validator{},
		&resource.Generator{},
	}
}
