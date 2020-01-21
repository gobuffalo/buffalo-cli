package refresh

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/develop"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
)

var _ plugins.Plugin = &Developer{}
var _ plugins.NamedCommand = &Developer{}
var _ develop.Developer = &Developer{}
