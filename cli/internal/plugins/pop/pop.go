package pop

import (
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/builder"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/built"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/generators/actions"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/generators/actiontest"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/generators/models"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/tester"
	"github.com/gobuffalo/buffalo-cli/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&actions.Generator{},
		&actiontest.Generator{},
		&builder.Builder{},
		&built.Initer{},
		&models.Generator{},
		&tester.Tester{},
	}
}
