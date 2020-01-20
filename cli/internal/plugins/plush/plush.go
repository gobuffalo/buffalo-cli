package plush

import (
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/plush/validator"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&resource.Generator{},
		&validator.Validator{},
	}
}
