package fizz

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/fizz/generate"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/fizz/setup"
	"github.com/gobuffalo/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&generate.Migration{},
		&setup.Setup{},
	}
}
