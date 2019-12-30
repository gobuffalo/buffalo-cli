package plugprint

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/internal/plugins"
)

type PluginScoper = plugins.PluginScoper
type Aliases = plugins.Aliases
type NamedCommand = plugins.NamedCommand

type SubCommander interface {
	SubCommands() []plugins.Plugin
}

// Describer is called by `Print` and can be
// implemented to print a short, single line,
// description of the plugin. `-h`
type Describer interface {
	Description() string
}

type FlagPrinter interface {
	PrintFlags(w io.Writer) error
}

type Hider interface {
	HidePlugin()
}

// UsagePrinter is called by `Print` and can be implemented
// to print a large block of usage information after the
// `Describer` interface is called. This is useful for printing
// flag information, links, and other messages to users.
type UsagePrinter interface {
	PrintUsage(w io.Writer) error
}
