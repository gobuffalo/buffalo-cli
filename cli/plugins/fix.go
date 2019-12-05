package plugins

import (
	"context"
	"fmt"
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
