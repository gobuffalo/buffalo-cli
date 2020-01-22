package fix

import (
	"context"
	"flag"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/spf13/pflag"
)

// Fixer is an optional interface a plugin can implement
// to be run with `buffalo fix`. This should update the application
// to the current version of the plugin.
// The expectation is fixing of only one major revision.
type Fixer interface {
	Fix(ctx context.Context, args []string) error
}

type Flagger interface {
	plugins.Plugin
	FixFlags() []*flag.Flag
}

type Pflagger interface {
	plugins.Plugin
	FixFlags() []*pflag.Flag
}
