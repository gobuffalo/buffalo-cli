package setup

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
)

func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	plugs := cmd.ScopedPlugins()
	if p := FindSetuperFromArgs(args, plugs); p != nil {
		return p.Setup(ctx, root, args[1:])
	}

	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return plugins.Wrap(cmd, err)
	}

	args = flags.Args()

	if cmd.help {
		return plugprint.Print(plugio.Stdout(plugs...), cmd)
	}

	err := cmd.run(ctx, root, args)
	return cmd.afterSetup(ctx, root, args, err)
}

func (cmd *Cmd) run(ctx context.Context, root string, args []string) error {
	if err := cmd.beforeSetup(ctx, root, args); err != nil {
		return plugins.Wrap(cmd, err)
	}

	for _, p := range cmd.ScopedPlugins() {
		if s, ok := p.(Setuper); ok {
			if err := s.Setup(ctx, root, args); err != nil {
				return plugins.Wrap(s, err)
			}
		}
	}

	return nil
}

func (cmd *Cmd) beforeSetup(ctx context.Context, root string, args []string) error {
	plugs := cmd.ScopedPlugins()

	for _, p := range plugs {
		if bb, ok := p.(BeforeSetuper); ok {
			if err := bb.BeforeSetup(ctx, root, args); err != nil {
				return plugins.Wrap(bb, err)
			}
		}
	}

	return nil
}

func (cmd *Cmd) afterSetup(ctx context.Context, root string, args []string, err error) error {
	plugs := cmd.ScopedPlugins()

	for _, p := range plugs {
		if bb, ok := p.(AfterSetuper); ok {
			if err := bb.AfterSetup(ctx, root, args, err); err != nil {
				return plugins.Wrap(bb, err)
			}
		}
	}
	return err
}
