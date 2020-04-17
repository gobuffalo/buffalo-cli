package cli

import (
	"os"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/clifix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/golang"
	"github.com/gobuffalo/buffalo-cli/v2/meta"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
)

var _ plugcmd.SubCommander = &Buffalo{}
var _ plugins.Plugin = &Buffalo{}
var _ plugins.Scoper = &Buffalo{}
var _ plugprint.Describer = &Buffalo{}

// Buffalo represents the `buffalo` cli.
type Buffalo struct {
	Plugins []plugins.Plugin
	root    string
}

// TODO move to the generated application code
// once packages are no longer internal
func insidePlugins(root string) []plugins.Plugin {
	var plugs []plugins.Plugin

	plugs = append(plugs, golang.Plugins()...)
	plugs = append(plugs, clifix.Plugins()...)
	return plugs
}

func NewFromRoot(root string) (*Buffalo, error) {
	b := &Buffalo{
		root: root,
	}

	b.Plugins = append(b.Plugins, cmds.AvailablePlugins(root)...)

	if meta.IsBuffalo(root) {
		b.Plugins = append(b.Plugins, insidePlugins(root)...)
	}

	// pre scope the plugins to thin the initial set
	plugs := b.ScopedPlugins()
	plugins.Sort(plugs)
	b.Plugins = plugs

	pfn := b.ScopedPlugins

	for _, p := range b.Plugins {
		switch t := p.(type) {
		case plugins.Needer:
			t.WithPlugins(pfn)
		}
	}

	return b, nil
}

func New() (*Buffalo, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return NewFromRoot(pwd)
}

func (b Buffalo) ScopedPlugins() []plugins.Plugin {
	root := b.root
	plugs := make([]plugins.Plugin, 0, len(b.Plugins))
	for _, p := range b.Plugins {
		switch t := p.(type) {
		case AvailabilityChecker:
			if !t.PluginAvailable(root) {
				continue
			}
		}
		plugs = append(plugs, p)
	}
	return plugs
}

func (b Buffalo) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range b.ScopedPlugins() {
		if _, ok := p.(Commander); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

// Name ...
func (Buffalo) PluginName() string {
	return "buffalo"
}

func (Buffalo) String() string {
	return "buffalo"
}

// Description ...
func (Buffalo) Description() string {
	return "Tools for working with Buffalo applications"
}
