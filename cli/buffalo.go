package cli

import (
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
