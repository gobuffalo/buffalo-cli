package cli

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/clifix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/bzr"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/develop"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/fix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/fizz"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/flect"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/generate"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/git"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/golang"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/grifts"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/i18n"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/info"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/mail"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/packr"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pkger"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/plush"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/refresh"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/soda"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/test"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/version"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
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

func NewWithInfo(inf here.Info) (*Buffalo, error) {
	root := inf.Dir
	b := &Buffalo{}

	pfn := func() []plugins.Plugin {
		return b.Plugins
	}

	b.Plugins = append(b.Plugins, build.Plugins()...)
	b.Plugins = append(b.Plugins, clifix.Plugins()...)
	b.Plugins = append(b.Plugins, develop.Plugins()...)
	b.Plugins = append(b.Plugins, fix.Plugins()...)
	b.Plugins = append(b.Plugins, fizz.Plugins()...)
	b.Plugins = append(b.Plugins, flect.Plugins()...)
	b.Plugins = append(b.Plugins, generate.Plugins()...)
	b.Plugins = append(b.Plugins, golang.Plugins()...)
	b.Plugins = append(b.Plugins, grifts.Plugins()...)
	b.Plugins = append(b.Plugins, i18n.Plugins()...)
	b.Plugins = append(b.Plugins, info.Plugins()...)
	b.Plugins = append(b.Plugins, mail.Plugins()...)
	b.Plugins = append(b.Plugins, packr.Plugins()...)
	b.Plugins = append(b.Plugins, pkger.Plugins()...)
	b.Plugins = append(b.Plugins, plush.Plugins()...)
	b.Plugins = append(b.Plugins, pop.Plugins()...)
	b.Plugins = append(b.Plugins, refresh.Plugins()...)
	b.Plugins = append(b.Plugins, resource.Plugins()...)
	b.Plugins = append(b.Plugins, soda.Plugins()...)
	b.Plugins = append(b.Plugins, test.Plugins()...)
	b.Plugins = append(b.Plugins, version.Plugins()...)

	if _, err := os.Stat(filepath.Join(root, "package.json")); err == nil {
		b.Plugins = append(b.Plugins, webpack.Plugins()...)
	}

	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		b.Plugins = append(b.Plugins, git.Plugins()...)
	}

	if _, err := os.Stat(filepath.Join(root, ".bzr")); err == nil {
		b.Plugins = append(b.Plugins, bzr.Plugins()...)
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
		f.WithHereInfo(inf)
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
