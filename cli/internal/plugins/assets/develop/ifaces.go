package develop

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/assets/internal/ifaces"
)

// Tooler returns the name of the asset tool to use.
// npm, yarnpkg, etc...
type Tooler = ifaces.Tooler

type Scripter = ifaces.Scripter
