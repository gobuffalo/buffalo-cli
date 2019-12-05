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

func (fc *fixCmd) Main(ctx context.Context, args []string) error {
	flags := cmdx.NewFlagSet("buffalo fix", cmdx.Stderr(ctx))
	flags.BoolVarP(&fix.YesToAll, "yes", "y", false, "update all without asking for confirmation")
	flags.BoolVarP(&fc.help, "help", "h", false, "print this help")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if fc.help {
		for _, p := range fc.Plugins {
			if _, ok := p.(plugins.Fixer); ok {
				fmt.Fprintf(fc.Stderr, "buffalo fix %s - [%s]\n", p.Name(), p)
			}
		}
		return nil
	}

	if len(args) > 0 {
		return fc.Plugins.Fix(ctx, args)
	}

	if err := fix.Run(); err != nil {
		return err
	}
	return fc.Plugins.Fix(ctx, args)
}
