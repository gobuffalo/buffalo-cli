package cli

import (
	"context"
	"time"

	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/info"
	"github.com/gobuffalo/clara/genny/rx"
	"github.com/gobuffalo/genny"
)

// Info implements the `buffalo info` command. Buffalo's checks
// are run first, then any plugins that implement plugins.Informer
// will be run in order at the end.
func (b *Buffalo) Info(ctx context.Context, args []string) error {
	var help bool
	flags := cmdx.NewFlagSet("buffalo info", cmdx.Stderr(ctx))
	flags.BoolVarP(&help, "help", "h", false, "print this help")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if help {
		flags.Usage()
		return nil
	}

	args = flags.Args()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	run := genny.WetRunner(ctx)

	out := cmdx.Stdout(ctx)

	opts := &rx.Options{
		Out: rx.NewWriter(out),
	}
	if err := run.WithNew(rx.New(opts)); err != nil {
		return err
	}

	iopts := &info.Options{
		Out: rx.NewWriter(out),
	}

	if err := run.WithNew(info.New(iopts)); err != nil {
		return err
	}

	if err := run.Run(); err != nil {
		return err
	}
	return b.Plugins.Info(ctx, args)
}
