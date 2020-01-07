package pop

import (
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/internal/actiongen"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/internal/modelgen"
	"github.com/gobuffalo/buffalo-cli/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Builder{},
		&Built{},
		&Cmd{},
		&Tester{},
		&actiongen.Generator{},
		&modelgen.Generator{},
	}
}
