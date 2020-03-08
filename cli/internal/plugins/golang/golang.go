package golang

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/golang/mod"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/golang/templates"
	"github.com/gobuffalo/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&mod.Initer{},
		&templates.Validator{},
	}
}
