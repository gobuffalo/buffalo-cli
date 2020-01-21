package develop

import (
	"context"
	"flag"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/spf13/pflag"
)

type Developer interface {
	// Develop will be called asyncronously with other implementations.
	// The context.Context should be listened to for cancellation.
	Develop(ctx context.Context, root string, args []string) error
}

type Flagger interface {
	plugins.Plugin
	DevelopFlags() []*flag.Flag
}

type Pflagger interface {
	plugins.Plugin
	DevelopFlags() []*pflag.Flag
}
