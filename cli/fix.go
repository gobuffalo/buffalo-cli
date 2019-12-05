package cli

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd/fix"
)

func (b *Buffalo) Fix(ctx context.Context, args []string) error {
	var help bool
	flags := cmdx.NewFlagSet("buffalo fix", cmdx.Stderr(ctx))
	flags.BoolVarP(&fix.YesToAll, "yes", "y", false, "update all without asking for confirmation")
	flags.BoolVarP(&help, "help", "h", false, "print this help")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if help {
		stderr := cmdx.Stderr(ctx)
		for _, p := range b.Plugins {
			if _, ok := p.(plugins.Fixer); ok {
				fmt.Fprintf(stderr, "buffalo fix %s - [%s]\n", p.Name(), p)
			}
		}
		return nil
	}

	if len(args) > 0 {
		return b.Plugins.Fix(ctx, args)
	}

	if err := fix.Run(); err != nil {
		return err
	}
	return b.Plugins.Fix(ctx, args)
}
