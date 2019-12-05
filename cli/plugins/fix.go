package plugins

import (
	"context"
	"fmt"
	"sort"

	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
)

// Fixer is an optional interface a plugin can implement
// to be run with `buffalo fix`. This should update the application
// to the current version of the plugin.
// The expectation is fixing of only one major revision.
type Fixer interface {
	Fix(ctx context.Context, args []string) error
}

// Fix runs any Fixers that are in the Plugins.
// If no arguments are provided it will run all fixers in the Plugins.
// Otherwise Fix will run the fixers for the arguments provided.
// 	buffalo fix
// 	buffalo fix plush pop
// 	buffalo fix -h
func (plugs Plugins) Fix(ctx context.Context, args []string) error {
	opts := struct {
		help bool
	}{}

	flags := cmdx.NewFlagSet("buffalo fix", cmdx.Stderr(ctx))
	flags.BoolVar(&opts.help, "h", false, "print this help")

	flags.Parse(args)

	args = flags.Args()

	// stderr := cmdx.Stderr(ctx)
	if opts.help {
		sort.Slice(plugs, func(i, j int) bool {
			return plugs[i].Name() < plugs[j].Name()
		})

		for _, p := range plugs {
			if _, ok := p.(Fixer); ok {
				// fmt.Fprintf(stderr, "%s %s - [%s]\n", flags.Name(), p.Name(), p)
			}
		}
		return nil
	}

	if len(args) > 0 {
		fixers := map[string]Fixer{}
		for _, p := range plugs {
			f, ok := p.(Fixer)
			if !ok {
				continue
			}

			fixers[p.Name()] = f
		}

		for _, a := range args {
			f, ok := fixers[a]
			if !ok {
				return fmt.Errorf("unknown fixer %s", a)
			}
			if err := f.Fix(ctx, []string{}); err != nil {
				return err
			}
		}
		return nil
	}

	for _, p := range plugs {
		f, ok := p.(Fixer)
		if !ok {
			continue
		}

		if err := f.Fix(ctx, args); err != nil {
			return err
		}
	}
	return nil
}
