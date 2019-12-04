package plugins

import (
	"context"
	"fmt"
	"sort"

	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
)

// Plugin is the most basic interface a plugin can implement.
type Plugin interface {
	// Name is the name of the plugin.
	// This will also be used for the cli sub-command
	// 	"pop" | "heroku" | "auth" | etc...
	Name() string
}

type Plugins []Plugin

// Len is the number of elements in the collection.
func (plugs Plugins) Len() int {
	return len(plugs)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (plugs Plugins) Less(i int, j int) bool {
	return plugs[i].Name() < plugs[j].Name()
}

// Swap swaps the elements with indexes i and j.
func (plugs Plugins) Swap(i int, j int) {
	plugs[i], plugs[j] = plugs[j], plugs[i]
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

// Generate will run the specified generator.
// 	buffalo generate -h
// 	buffalo generate pop ...
func (plugs Plugins) Generate(ctx context.Context, args []string) error {
	opts := struct {
		help bool
	}{}

	flags := cmdx.NewFlagSet("buffalo generate", cmdx.Stderr(ctx))
	flags.BoolVar(&opts.help, "h", false, "print this help")
	flags.Parse(args)

	args = flags.Args()
	if opts.help || len(args) == 0 {
		sort.Slice(plugs, func(i, j int) bool {
			return plugs[i].Name() < plugs[j].Name()
		})

		// stderr := cmdx.Stderr(ctx)
		for _, p := range plugs {
			if _, ok := p.(Generator); ok {
				// fmt.Fprintf(stderr, "%s %s - [%s]\n", flags.Name(), p.Name(), p)
			}
		}
		return nil
	}

	arg := args[0]
	if len(args) > 0 {
		args = args[1:]
	}

	for _, p := range plugs {
		f, ok := p.(Generator)
		if !ok {
			continue
		}
		if p.Name() != arg {
			continue
		}

		return f.Generate(ctx, args)
	}
	return fmt.Errorf("unknown generator %s", arg)
}
