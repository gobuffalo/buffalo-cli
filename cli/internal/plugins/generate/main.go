package generate

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
)

// Main implements cli.Cmd and is the entry point for `buffalo generate`
func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	ioe := plugins.CtxIO(ctx)
	if len(args) == 0 {
		if err := plugprint.Print(ioe.Stdout(), cmd); err != nil {
			return err
		}
		return fmt.Errorf("no command provided")
	}

	if len(args) == 1 && args[0] == "-h" {
		return plugprint.Print(ioe.Stdout(), cmd)
	}

	plugs := cmd.ScopedPlugins()

	n := args[0]
	cmds := plugins.Commands(plugs)
	p, err := cmds.Find(n)
	if err != nil {
		return err
	}

	b, ok := p.(Generator)
	if !ok {
		return fmt.Errorf("unknown command %q", n)
	}

	return b.Generate(ctx, root, args[1:])
}
