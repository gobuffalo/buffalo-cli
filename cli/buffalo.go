package cli

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/assets"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/bzr"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/fixcmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/git"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/golang"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/infocmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/packr"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pkger"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/plush"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/versioncmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here/there"
)

var _ plugins.Plugin = &Buffalo{}
var _ plugprint.SubCommander = &Buffalo{}
var _ plugprint.Describer = &Buffalo{}
var _ plugprint.WithPlugins = &Buffalo{}

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

	info, err := there.Current()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filepath.Join(info.Root, ".git")); err == nil {
		b.Plugins = append(b.Plugins,
			&git.Buffalo{
				IO: b,
			})
	}
	if _, err := os.Stat(filepath.Join(info.Root, ".bzr")); err == nil {
		b.Plugins = append(b.Plugins,
			&bzr.Buffalo{
				IO: b,
			})
	}
	return b, nil
}

func (b Buffalo) WithPlugins() plugins.Plugins {
	return b.Plugins
}

func (b Buffalo) SubCommands() plugins.Plugins {
	var plugs plugins.Plugins
	for _, p := range b.WithPlugins() {
		if _, ok := p.(Command); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
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
