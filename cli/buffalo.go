package cli

import (
	"context"
	"io"
	"os"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd"
	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd/fix"
	"github.com/spf13/pflag"
)

// Buffalo represents the `buffalo` cli.
type Buffalo struct {
	context.Context
	flags   *pflag.FlagSet
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Plugins plugins.Plugins
	version bool
	help    bool
}

func New(ctx context.Context) (*Buffalo, error) {
	b := &Buffalo{
		Context: ctx,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
	b.setFlags()
	return b, nil
}

func (b *Buffalo) Flags() *pflag.FlagSet {
	if b.flags == nil {
		b.setFlags()
	}
	return b.flags
}

func (b *Buffalo) setFlags() {
}

func (b *Buffalo) Fix(ctx context.Context, args []string) error {
	flags := cmdx.NewFlagSet("buffalo fix", cmdx.Stderr(ctx))
	flags.BoolVar(&fix.YesToAll, "y", false, "update all without asking for confirmation")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if err := fix.Run(); err != nil {
		return err
	}
	return b.Plugins.Fix(ctx, args)
}

func (b *Buffalo) Main(ctx context.Context, args []string) error {
	// flags := cmdx.NewFlagSet("buffalo", cmdx.Stderr(ctx))
	// flags.BoolVar(&b.version, "v", false, "display version")
	// flags.BoolVar(&b.help, "h", false, "display help")
	if len(args) > 0 {
		switch args[0] {
		case "fix":
			return b.Fix(ctx, args[1:])
		case "info":
			return b.Info(ctx, args[1:])
		}
	}

	c := cmd.RootCmd
	c.SetArgs(args)
	return c.Execute()
}
