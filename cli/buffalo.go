package cli

import (
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/assets"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/fixcmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/git"
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
	plugins.IO
	Plugins plugins.Plugins
}

func New() (*Buffalo, error) {
	b := &Buffalo{
		IO: plugins.NewIO(),
	}

	pfn := func() plugins.Plugins {
		return b.Plugins
	}
	b.Plugins = append(b.Plugins,
		&assets.Builder{
			IO: b,
		},
		&buildcmd.BuildCmd{
			IO:      b,
			Parent:  b,
			Plugins: pfn,
		},
		&fixcmd.FixCmd{
			IO:      b,
			Parent:  b,
			Plugins: pfn,
		},
		&infocmd.InfoCmd{
			IO:      b,
			Parent:  b,
			Plugins: pfn,
		},
		&versioncmd.VersionCmd{
			IO:     b,
			Parent: b,
		},
		&plush.Buffalo{},
		&golang.Templates{},
		&packr.Buffalo{},
		&pkger.Buffalo{},
	)

	if _, err := exec.LookPath("git"); err == nil {
		b.Plugins = append(b.Plugins,
			&git.Buffalo{
				IO: b,
			})
	}
	return b, nil
}

func (b Buffalo) WithPlugins() plugins.Plugins {
	return b.Plugins
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
