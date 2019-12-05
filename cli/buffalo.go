package cli

import (
	"context"
	"io"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
)

// Buffalo represents the `buffalo` cli.
type Buffalo struct {
	context.Context
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Plugins plugins.Plugins
	help    bool
}

func New(ctx context.Context) (*Buffalo, error) {
	b := &Buffalo{
		Context: ctx,
		Stdin:   cmdx.Stdin(ctx),
		Stdout:  cmdx.Stdout(ctx),
		Stderr:  cmdx.Stderr(ctx),
	}
	b.Plugins = append(b.Plugins,
		&versionCmd{Buffalo: b},
		&fixCmd{Buffalo: b},
		&infoCmd{Buffalo: b},
	)
	return b, nil
}

func (Buffalo) Name() string {
	return "buffalo"
}

func (Buffalo) Description() string {
	return "Tools for working with Buffalo applications"
}
