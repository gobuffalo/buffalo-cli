package cli

import (
	"context"
	"io"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
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
		Stdin:   cmdx.Stdin(ctx),
		Stdout:  cmdx.Stdout(ctx),
		Stderr:  cmdx.Stderr(ctx),
	}
	b.Plugins = append(b.Plugins,
		&fixCmd{Buffalo: b},
		&infoCmd{Buffalo: b},
		&versionCmd{Buffalo: b},
	)
	return b, nil
}

func (b *Buffalo) Main(ctx context.Context, args []string) error {
	c := cmd.RootCmd
	if len(args) == 0 {
		c.SetArgs(args)
		return c.Execute()
	}

	cmds := b.Plugins.Commands()
	if c, err := cmds.Find(args[0]); err == nil {
		return c.Main(ctx, args[1:])
	}

	c.SetArgs(args)
	return c.Execute()
}
