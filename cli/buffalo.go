package cli

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/assets"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/bzr"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/fixcmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/fizz"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/flect"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/generatecmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/git"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/golang"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/grifts"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/i18n"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/infocmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/mail"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pkger"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/plush"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/soda"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/testcmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/versioncmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
)

var _ plugins.Plugin = &Buffalo{}
var _ plugins.PluginScoper = &Buffalo{}
var _ plugprint.Describer = &Buffalo{}
var _ plugprint.SubCommander = &Buffalo{}

// Buffalo represents the `buffalo` cli.
type Buffalo struct {
	plugins.Plugins
}

func NewWithInfo(info here.Info) (*Buffalo, error) {
	b := &Buffalo{}

	pfn := func() []plugins.Plugin {
		return b.Plugins
	}
	b.Plugins = append(b.Plugins,
		&assets.Builder{},
		&buildcmd.BuildCmd{},
		&buildcmd.MainFile{},
		&fixcmd.FixCmd{},
		&fizz.MigrationGen{},
		&flect.Buffalo{},
		&generatecmd.Command{},
		&golang.Templates{},
		&i18n.Generator{},
		&infocmd.InfoCmd{},
		&mail.Generator{},
		&pkger.Buffalo{},
		&resource.Generator{},
		&testcmd.TestCmd{},
		&versioncmd.VersionCmd{},
		// &packr.Buffalo{},
	)
	b.Plugins = append(b.Plugins, grifts.Plugins()...)
	b.Plugins = append(b.Plugins, plush.Plugins()...)
	b.Plugins = append(b.Plugins, pop.Plugins()...)
	b.Plugins = append(b.Plugins, soda.Plugins()...)

	if _, err := os.Stat(filepath.Join(info.Dir, ".git")); err == nil {
		b.Plugins = append(b.Plugins, git.Plugins()...)
	}
	if _, err := os.Stat(filepath.Join(info.Dir, ".bzr")); err == nil {
		b.Plugins = append(b.Plugins, &bzr.Buffalo{})
	}

	sort.Slice(b.Plugins, func(i, j int) bool {
		return b.Plugins[i].Name() < b.Plugins[j].Name()
	})

	pfn = func() []plugins.Plugin {
		return b.Plugins
	}

	for _, b := range b.Plugins {
		f, ok := b.(plugins.PluginNeeder)
		if !ok {
			continue
		}
		f.WithPlugins(pfn)
	}

	for _, b := range b.Plugins {
		f, ok := b.(WithHere)
		if !ok {
			continue
		}
		f.WithHereInfo(info)
	}

	return b, nil
}

func New() (*Buffalo, error) {
	info, err := here.Current()
	if err != nil {
		return nil, err
	}
	return NewWithInfo(info)
}

func (b Buffalo) ScopedPlugins() []plugins.Plugin {
	return b.Plugins
}

func (b Buffalo) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range b.ScopedPlugins() {
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
