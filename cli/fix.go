package cli

import (
	"context"

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
		flags.Usage()
		return nil
	}

	if err := fix.Run(); err != nil {
		return err
	}
	return b.Plugins.Fix(ctx, args)
}
