package build

import (
	"context"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
)

func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	plugs := cmd.ScopedPlugins()
	if sub := FindBuilderFromArgs(args, plugs); sub != nil {
		return sub.Build(ctx, root, args[1:])
	}

	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return plugins.Wrap(cmd, err)
	}

	if cmd.help {
		return plugprint.Print(plugio.Stdout(cmd.ScopedPlugins()...), cmd)
	}

	err := cmd.run(ctx, root, args)
	return cmd.afterBuild(ctx, root, args, err)
}

func (cmd *Cmd) run(ctx context.Context, root string, args []string) error {
	info, err := here.Dir(root)
	if err != nil {
		return plugins.Wrap(cmd, err)
	}

	if err = cmd.beforeBuild(ctx, root, args); err != nil {
		return plugins.Wrap(cmd, err)
	}

	plugs := cmd.ScopedPlugins()

	if err := cmd.pack(ctx, info, plugs); err != nil {
		return plugins.Wrap(cmd, err)
	}

	return cmd.build(ctx, root, args)
}

func (cmd *Cmd) beforeBuild(ctx context.Context, root string, args []string) error {
	plugs := cmd.ScopedPlugins()
	for _, p := range plugs {
		if bb, ok := p.(BeforeBuilder); ok {
			if err := bb.BeforeBuild(ctx, root, args); err != nil {
				return plugins.Wrap(p, err)
			}
		}
	}
	return nil
}

func (cmd *Cmd) afterBuild(ctx context.Context, root string, args []string, err error) error {
	plugs := cmd.ScopedPlugins()
	for _, p := range plugs {
		if bb, ok := p.(AfterBuilder); ok {
			if err := bb.AfterBuild(ctx, root, args, err); err != nil {
				return plugins.Wrap(p, err)
			}
		}
	}
	return err
}
