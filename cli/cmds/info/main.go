package info

import (
	"context"

	"github.com/gobuffalo/plugins"
)

// Main implements the `buffalo info` command. Buffalo's checks
// are run first, then any plugins that implement plugins.Informer
// will be run in order at the end.
func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	plugs := cmd.ScopedPlugins()

	if p := FindInformerFromArgs(args, plugs); p != nil {
		return p.Info(ctx, root, args[1:])
	}

	for _, p := range plugs {
		if i, ok := p.(Informer); ok {
			if err := i.Info(ctx, root, args); err != nil {
				return plugins.Wrap(i, err)
			}
		}
	}
	return nil
}
