package fix

import (
	"context"
	"flag"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/spf13/pflag"
)

// Fixer is an optional interface a plugin can implement
// to be run with `buffalo fix`. This should update the application
// to the current version of the plugin.
// The expectation is fixing of only one major revision.
type Fixer interface {
	plugins.Plugin
	Fix(ctx context.Context, root string, args []string) error
}

type Flagger interface {
	plugins.Plugin
	FixFlags() []*flag.Flag
}

type Pflagger interface {
	plugins.Plugin
	FixFlags() []*pflag.Flag
}

type Namer interface {
	Fixer
	CmdName() string
}

type Aliaser interface {
	Fixer
	CmdAliases() []string
}

type Stdouter = plugio.Outer
