package cli

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/assets"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/bzr"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/fixcmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/flect"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/git"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/golang"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/grifts"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/infocmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pkger"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/plush"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/versioncmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
)

var _ plugins.Plugin = &Buffalo{}
var _ plugprint.SubCommander = &Buffalo{}
var _ plugprint.Describer = &Buffalo{}
var _ plugprint.Plugins = &Buffalo{}

// Buffalo represents the `buffalo` cli.
type Buffalo struct {
	plugins.Plugins
}

func New() (*Buffalo, error) {
	b := &Buffalo{}

	pfn := func() []plugins.Plugin {
		return b.Plugins
	}
	b.Plugins = append(b.Plugins,
		&flect.Flect{},
		&pop.Buffalo{},
		&grifts.Buffalo{},
		&assets.Builder{},
		&buildcmd.BuildCmd{
			PluginsFn: pfn,
		},
		&fixcmd.FixCmd{
			Parent:    b,
			PluginsFn: pfn,
		},
		&infocmd.InfoCmd{
			Parent:    b,
			PluginsFn: pfn,
		},
		&versioncmd.VersionCmd{
			Parent: b,
		},
		&plush.Buffalo{},
		&golang.Templates{},
		// &packr.Buffalo{},
		&pkger.Buffalo{
			PluginsFn: pfn,
		},
	)

	info, err := here.Current()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filepath.Join(info.Root, ".git")); err == nil {
		b.Plugins = append(b.Plugins, &git.Buffalo{})
	}
	if _, err := os.Stat(filepath.Join(info.Root, ".bzr")); err == nil {
		b.Plugins = append(b.Plugins, &bzr.Buffalo{})
	}
	return b, nil
}

func (b Buffalo) WithPlugins() []plugins.Plugin {
	return b.Plugins
}

func (b Buffalo) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
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
