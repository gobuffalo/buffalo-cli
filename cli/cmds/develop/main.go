package develop

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"golang.org/x/sync/errgroup"
)

func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	plugs := cmd.ScopedPlugins()
	if sub := FindDeveloperFromArgs(args, plugs); sub != nil {
		return sub.Develop(ctx, root, args[1:])
	}

	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return plugins.Wrap(cmd, err)
	}

	if cmd.help {
		return plugprint.Print(plugio.Stdout(cmd.ScopedPlugins()...), cmd)
	}

	args = flags.Args()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &errgroup.Group{}

	for _, p := range cmd.ScopedPlugins() {
		if d, ok := p.(Developer); ok {
			wg.Go(func() error {
				return d.Develop(ctx, root, args)
			})
		}
	}

	return wg.Wait()
}
