package built

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/plugins"
)

// Initer is invoked in when an application binary
// built with `buffalo build` is executed. This hook
// is executed before any flags are parsed or sub-commands
// are run.
type Initer interface {
	BuiltInit(ctx context.Context, args []string) error
}

type IniterFn func(ctx context.Context, args []string) error

func (i IniterFn) BuiltInit(ctx context.Context, args []string) error {
	return i(ctx, args)
}

func WithIniter(p plugins.Plugin, fn IniterFn) plugins.Plugin {
	type wi struct {
		IniterFn
		plugins.Plugin
	}
	return wi{
		Plugin:   p,
		IniterFn: fn,
	}
}
