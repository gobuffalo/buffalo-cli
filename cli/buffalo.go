package cli

import (
	"io"
	"os"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/assets"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/fixcmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/golang"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/infocmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/packr"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pkger"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/plush"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/versioncmd"
	"github.com/gobuffalo/buffalo-cli/cli/plugins"
)

// Buffalo represents the `buffalo` cli.
type Buffalo struct {
	Plugins plugins.Plugins
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
}

func New() (*Buffalo, error) {
	b := &Buffalo{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	pfn := func() plugins.Plugins {
		return b.Plugins
	}
	b.Plugins = append(b.Plugins,
		&assets.Builder{},
		&buildcmd.BuildCmd{
			Parent:  b,
			Plugins: pfn,
		},
		&fixcmd.FixCmd{
			Parent:  b,
			Plugins: pfn,
		},
		&infocmd.InfoCmd{
			Parent:  b,
			Plugins: pfn,
		},
		&versioncmd.VersionCmd{
			Parent: b,
		},
		&plush.Buffalo{},
		&golang.Templates{},
		&packr.Buffalo{},
		&pkger.Buffalo{},
	)
	return b, nil
}

func (b *Buffalo) Stdin() io.Reader {
	if b.stdin == nil {
		return os.Stdin
	}
	return b.stdin
}

func (b *Buffalo) Stdout() io.Writer {
	if b.stdout == nil {
		return os.Stdout
	}
	return b.stdout
}

func (b *Buffalo) Stderr() io.Writer {
	if b.stderr == nil {
		return os.Stderr
	}
	return b.stderr
}

func (b *Buffalo) SetStdin(r io.Reader) {
	b.stdin = r
}

func (b *Buffalo) SetStdout(w io.Writer) {
	b.stdout = w
}

func (b *Buffalo) SetStderr(w io.Writer) {
	b.stderr = w
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
