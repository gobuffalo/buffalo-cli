package assets

import (
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/assets/builder"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/assets/developer"
	"github.com/gobuffalo/buffalo-cli/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&builder.Builder{},
		&developer.Developer{},
	}
}
