package generate

import (
	"context"
	"fmt"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugfind"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
)

// Main implements cli.Cmd and is the entry point for `buffalo generate`
func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	stdout := plugio.Stdout(cmd.ScopedPlugins()...)
	if len(args) == 0 {
		if err := plugprint.Print(stdout, cmd); err != nil {
			return err
		}
		return fmt.Errorf("no command provided")
	}

	if len(args) == 1 && args[0] == "-h" {
		return plugprint.Print(stdout, cmd)
	}

	plugs := cmd.ScopedPlugins()

	n := args[0]

	fn := plugfind.Background()
	fn = byGenerator(fn)
	fn = plugcmd.ByNamer(fn)
	fn = plugcmd.ByAliaser(fn)

	p := fn.Find(n, plugs)
	if p == nil {
		return fmt.Errorf("unknown command %q", n)
	}

	b, ok := p.(Generator)
	if !ok {
		return fmt.Errorf("unknown command %q", n)
	}

	return b.Generate(ctx, root, args[1:])
}

func byGenerator(f plugfind.Finder) plugfind.Finder {
	fn := func(name string, plugs []plugins.Plugin) plugins.Plugin {
		p := f.Find(name, plugs)
		if p == nil {
			return nil
		}
		if c, ok := p.(Generator); ok {
			if c.PluginName() == name {
				return p
			}
		}
		return nil
	}
	return plugfind.FinderFn(fn)
}
