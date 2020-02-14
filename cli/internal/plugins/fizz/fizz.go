package fizz

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/fizz/setup"
	"github.com/gobuffalo/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&MigrationGen{},
		&setup.Setup{},
	}
}
