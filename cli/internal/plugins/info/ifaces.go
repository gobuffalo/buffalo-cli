package info

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/plugins"
)

// Informer can be implemented to add more checks
// to `buffalo info`
type Informer interface {
	plugins.Plugin
	Info(ctx context.Context, args []string) error
}
