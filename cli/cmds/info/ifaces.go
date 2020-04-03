package info

import (
	"context"

	"github.com/gobuffalo/plugins"
)

// Informer can be implemented to add more checks
// to `buffalo info`
type Informer interface {
	plugins.Plugin
	Info(ctx context.Context, root string, args []string) error
}

type Namer interface {
	Informer
	CmdName() string
}

type Aliaser interface {
	Informer
	CmdAliases() []string
}
