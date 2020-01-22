package plush

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/plush/validator"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&resource.Generator{},
		&validator.Validator{},
	}
}
