package plush

import (
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/plush/internal/resourcegen"
	"github.com/gobuffalo/buffalo-cli/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Validator{},
		&resourcegen.Generator{},
	}
}
