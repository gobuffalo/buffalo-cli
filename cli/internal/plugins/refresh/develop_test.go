package refresh

import (
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/develop"
	"github.com/gobuffalo/buffalo-cli/plugins"
)

var _ plugins.Plugin = &Developer{}
var _ plugins.NamedCommand = &Developer{}
var _ develop.Developer = &Developer{}
