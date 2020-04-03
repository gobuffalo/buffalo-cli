package fix

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
)

func (fc *Cmd) Main(ctx context.Context, root string, args []string) error {
	plugs := fc.ScopedPlugins()

	if p := FindFixerFromArgs(args, plugs); p != nil {
		return p.Fix(ctx, root, args[1:])
	}

	flags := fc.Flags()

	if err := flags.Parse(args); err != nil {
		return plugins.Wrap(fc, err)
	}

	if fc.help {
		return plugprint.Print(plugio.Stdout(fc.ScopedPlugins()...), fc)
	}

	return fc.run(ctx, root, args)
}

// Fix runs any Fixers that are in the Plugins.
// If no arguments are provided it will run all fixers in the Plugins.
// Otherwise Fix will run the fixers for the arguments provided.
// 	buffalo fix
// 	buffalo fix plush pop
// 	buffalo fix -h
func (fc *Cmd) run(ctx context.Context, root string, args []string) error {
	plugs := fc.ScopedPlugins()

	for _, p := range plugs {
		f, ok := p.(Fixer)
		if !ok {
			continue
		}

		if err := f.Fix(ctx, root, []string{}); err != nil {
			return plugins.Wrap(f, err)
		}
	}

	return nil
}
