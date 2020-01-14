package pop

import (
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/internal/generators/actions"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/internal/generators/actiontest"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/internal/generators/models"
	"github.com/gobuffalo/buffalo-cli/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Builder{},
		&Built{},
		&Tester{},
		&actions.Generator{},
		&actiontest.Generator{},
		&models.Generator{},
	}
}
