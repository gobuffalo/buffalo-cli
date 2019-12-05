package cli

import (
	"context"
	"io"
	"os"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd"
)

// Buffalo represents the `buffalo` cli.
type Buffalo struct {
	context.Context
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Plugins plugins.Plugins
}

func New(ctx context.Context) (*Buffalo, error) {
	b := &Buffalo{
		Context: ctx,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
	return b, nil
}

func (b *Buffalo) Main(ctx context.Context, args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "fix":
			return b.Fix(ctx, args[1:])
		case "info":
			return b.Info(ctx, args[1:])
		case "version":
			return b.Version(ctx, args[1:])
		}
	}
	c := cmd.RootCmd
	c.SetArgs(args)
	return c.Execute()
}
