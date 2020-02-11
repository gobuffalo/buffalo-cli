package develop

import (
	"context"
	"fmt"
	"strings"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugfind"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"golang.org/x/sync/errgroup"
)

func (cmd *Cmd) SubCommand(ctx context.Context, root string, name string, args []string) error {
	plugs := cmd.SubCommands()

	fn := plugfind.Background()
	fn = byDeveloper(fn)
	fn = plugcmd.ByNamer(fn)
	fn = plugcmd.ByAliaser(fn)

	p := fn.Find(name, plugs)
	if p == nil {
		return fmt.Errorf("%s is not a developer", name)
	}

	d, ok := p.(Developer)
	if !ok {
		return fmt.Errorf("%s is not a developer", name)
	}

	return d.Develop(ctx, root, args)
}

func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	if len(args) == 1 && args[0] == "-h" {
		return plugprint.Print(plugio.Stdout(cmd.ScopedPlugins()...), cmd)
	}

	if len(args) > 0 {
		for _, n := range args {
			if strings.HasPrefix(n, "-") {
				continue
			}
			return cmd.SubCommand(ctx, root, n, args[1:])
		}

	}

	flags := cmd.Flags()
	if err := flags.Parse(args); err != nil {
		return err
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

func byDeveloper(f plugfind.Finder) plugfind.Finder {
	fn := func(name string, plugs []plugins.Plugin) plugins.Plugin {
		p := f.Find(name, plugs)
		if p == nil {
			return nil
		}
		if c, ok := p.(Developer); ok {
			if c.PluginName() == name {
				return p
			}
		}
		return nil
	}
	return plugfind.FinderFn(fn)
}
