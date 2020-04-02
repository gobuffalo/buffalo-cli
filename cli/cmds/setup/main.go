package setup

import (
	"context"
	"fmt"
	"strings"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugfind"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/markbates/safe"
)

func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	if len(args) > 0 {
		if !strings.HasPrefix(args[0], "-") {
			return cmd.SubCommand(ctx, root, args[0], args)
		}
	}

	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return plugins.Wrap(cmd, err)
	}

	args = flags.Args()

	if cmd.help {
		return plugprint.Print(plugio.Stdout(cmd.ScopedPlugins()...), cmd)
	}

	if err := cmd.beforeSetup(ctx, root, args); err != nil {
		return cmd.afterSetup(ctx, root, args, err)
	}

	for _, p := range cmd.ScopedPlugins() {
		if s, ok := p.(Setuper); ok {
			if err := s.Setup(ctx, root, args); err != nil {
				return cmd.afterSetup(ctx, root, args, err)
			}
		}
	}

	return cmd.afterSetup(ctx, root, args, nil)
}

func (cmd *Cmd) beforeSetup(ctx context.Context, root string, args []string) error {
	plugs := cmd.ScopedPlugins()
	for _, p := range plugs {
		if bb, ok := p.(BeforeSetuper); ok {
			err := safe.RunE(func() error {
				return bb.BeforeSetup(ctx, root, args)
			})
			if err != nil {
				return plugins.Wrap(cmd, err)
			}
		}
	}
	return nil
}

func (cmd *Cmd) afterSetup(ctx context.Context, root string, args []string, err error) error {
	plugs := cmd.ScopedPlugins()
	for _, p := range plugs {
		if bb, ok := p.(AfterSetuper); ok {
			err := safe.RunE(func() error {
				return bb.AfterSetup(ctx, root, args, err)
			})
			if err != nil {
				return plugins.Wrap(cmd, err)
			}
		}
	}
	return err
}

func bySetuper(f plugfind.Finder) plugfind.Finder {
	fn := func(name string, plugs []plugins.Plugin) plugins.Plugin {
		p := f.Find(name, plugs)
		if p == nil {
			return nil
		}
		if c, ok := p.(Setuper); ok {
			if c.PluginName() == name {
				return p
			}
		}
		return nil
	}
	return plugfind.FinderFn(fn)
}

func (cmd *Cmd) SubCommand(ctx context.Context, root string, name string, args []string) error {
	plugs := cmd.SubCommands()

	fn := plugfind.Background()
	fn = bySetuper(fn)
	fn = plugcmd.ByNamer(fn)
	fn = plugcmd.ByAliaser(fn)

	p := fn.Find(name, plugs)

	d, ok := p.(Setuper)
	if !ok {
		return fmt.Errorf("%s unknown command", name)
	}

	return d.Setup(ctx, root, args[1:])
}
