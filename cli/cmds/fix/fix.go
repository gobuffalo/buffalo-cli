package fix

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/golang/mainfix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/refresh"
	"github.com/gobuffalo/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Cmd{},
		&mainfix.Cmd{},
		&refresh.Fixer{},
	}
}
