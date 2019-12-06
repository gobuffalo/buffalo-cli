package cli

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd/fix"
)

type fixCmd struct {
	*Buffalo
	help bool
}

func (fc *fixCmd) Name() string {
	return "fix"
}

func (fc *fixCmd) Description() string {
	return "Attempt to fix a Buffalo application's API to match version in go.mod"
}

// Fix runs any Fixers that are in the Plugins.
// If no arguments are provided it will run all fixers in the Plugins.
// Otherwise Fix will run the fixers for the arguments provided.
// 	buffalo fix
// 	buffalo fix plush pop
// 	buffalo fix -h
func (fc *fixCmd) plugins(ctx context.Context, args []string) error {
	plugs := fc.Plugins
	if len(args) > 0 {
		fixers := map[string]plugins.Fixer{}
		for _, p := range plugs {
			f, ok := p.(plugins.Fixer)
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
		f, ok := p.(plugins.Fixer)
		if !ok {
			continue
		}

		if err := f.Fix(ctx, args); err != nil {
			return err
		}
	}
	return nil
}

func (fc *fixCmd) Main(ctx context.Context, args []string) error {
	flags := cmdx.NewFlagSet("buffalo fix", fc.Stderr)
	flags.BoolVarP(&fix.YesToAll, "yes", "y", false, "update all without asking for confirmation")
	flags.BoolVarP(&fc.help, "help", "h", false, "print this help")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if fc.help {
		var plugs plugins.Plugins
		for _, p := range fc.Plugins {
			if _, ok := p.(plugins.Fixer); ok {
				plugs = append(plugs, p)
			}
		}
		return cmdx.Print(fc.Stderr, fc.Buffalo.Name(), fc, plugs, flags)
	}

	if len(args) > 0 {
		return fc.plugins(ctx, args)
	}

	if err := fix.Run(); err != nil {
		return err
	}
	return fc.plugins(ctx, args)
}
