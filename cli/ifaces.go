package cli

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"
)

type Aliaser = plugcmd.Aliaser
type Commander = plugcmd.Commander
type Needer = plugins.Needer
type StderrNeeder = plugio.ErrNeeder
type StdinNeeder = plugio.InNeeder
type StdoutNeeder = plugio.OutNeeder

// AvailabilityChecker
type AvailabilityChecker interface {
	PluginAvailable(root string) bool
}
