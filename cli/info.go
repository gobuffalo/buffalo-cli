package cli

import (
	"context"
	"time"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/info"
	"github.com/gobuffalo/clara/genny/rx"
	"github.com/gobuffalo/genny"
)

type infoCmd struct {
	*Buffalo
	help bool
}

func (ic *infoCmd) Name() string {
	return "info"
}

func (ic *infoCmd) Description() string {
	return "Print diagnostic information (useful for debugging)"
}

// Info runs all of the plugins that implement the
// `Informer` interface in order.
func (ic *infoCmd) plugins(ctx context.Context, args []string) error {
	plugs := ic.Plugins
	for _, p := range plugs {
		i, ok := p.(plugins.Informer)
		if !ok {
			continue
		}
		if err := i.Info(ctx, args); err != nil {
			return err
		}
	}
	return nil
}

// Main implements the `buffalo info` command. Buffalo's checks
// are run first, then any plugins that implement plugins.Informer
// will be run in order at the end.
func (ic *infoCmd) Main(ctx context.Context, args []string) error {
	flags := cmdx.NewFlagSet("buffalo info", cmdx.Stderr(ctx))
	flags.BoolVarP(&ic.help, "help", "h", false, "print this help")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if ic.help {
		return cmdx.Print(ic.Stderr, ic.Buffalo.Name(), ic, nil, flags)
	}

	args = flags.Args()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	run := genny.WetRunner(ctx)

	out := ic.Stdout

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
	return ic.plugins(ctx, args)
}
