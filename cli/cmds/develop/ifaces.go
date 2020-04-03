package develop

import (
	"context"
	"flag"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/spf13/pflag"
)

type Developer interface {
	plugins.Plugin
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

type Namer interface {
	Developer
	CmdName() string
}

type Aliaser interface {
	Developer
	CmdAliases() []string
}

type Stdouter = plugio.Outer
