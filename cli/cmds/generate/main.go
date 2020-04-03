package generate

import (
	"context"

	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
)

// Main implements cli.Cmd and is the entry point for `buffalo generate`
func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	plugs := cmd.ScopedPlugins()

	if p := FindGeneratorFromArgs(args, plugs); p != nil {
		return p.Generate(ctx, root, args[1:])
	}

	stdout := plugio.Stdout(cmd.ScopedPlugins()...)
	return plugprint.Print(stdout, cmd)
}
