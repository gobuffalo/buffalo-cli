package cli

import (
	"io"
	"os"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/assets"
	"github.com/gobuffalo/buffalo-cli/cli/plugins"
)

// Buffalo represents the `buffalo` cli.
type Buffalo struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Plugins plugins.Plugins
}

func New() (*Buffalo, error) {
	b := &Buffalo{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	b.Plugins = append(b.Plugins,
		&versionCmd{Buffalo: b},
		&fixCmd{Buffalo: b},
		&buildCmd{Buffalo: b},
		&assets.Assets{
			Stdin:  b.Stdin,
			Stdout: b.Stdout,
			Stderr: b.Stderr,
		},
		&infoCmd{Buffalo: b},
	)
	return b, nil
}

// Name ...
func (Buffalo) Name() string {
	return "buffalo"
}

func (Buffalo) String() string {
	return "buffalo"
}

// Description ...
func (Buffalo) Description() string {
	return "Tools for working with Buffalo applications"
}
