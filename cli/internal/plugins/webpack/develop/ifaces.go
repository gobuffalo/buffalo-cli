package develop

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack/internal/ifaces"
	"github.com/gobuffalo/plugins/plugio"
)

// Tooler returns the name of the asset tool to use.
// npm, yarnpkg, etc...
type Tooler = ifaces.Tooler

type Scripter = ifaces.Scripter

type Stdouter = plugio.Outer
type Stderrer = plugio.Errer
type Stdiner = plugio.Inner
