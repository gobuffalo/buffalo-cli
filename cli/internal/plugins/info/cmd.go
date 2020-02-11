package info

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &Cmd{}
var _ plugins.Needer = &Cmd{}
var _ plugins.Scoper = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}

type Cmd struct {
	flags     *pflag.FlagSet
	pluginsFn plugins.Feeder
	help      bool
}

func (cmd *Cmd) WithPlugins(f plugins.Feeder) {
	cmd.pluginsFn = f
}

func (cmd *Cmd) PluginName() string {
	return "info"
}

func (cmd *Cmd) Description() string {
	return "Print diagnostic information (useful for debugging)"
}

// Info runs all of the plugins that implement the
// `Informer` interface in order.
func (cmd *Cmd) plugins(ctx context.Context, root string, args []string) error {
	for _, p := range cmd.ScopedPlugins() {
		i, ok := p.(Informer)
		if !ok {
			continue
		}
		if err := i.Info(ctx, root, args); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin

	if cmd.pluginsFn == nil {
		return plugs
	}
	for _, p := range cmd.pluginsFn() {
		if i, ok := p.(Informer); ok {
			plugs = append(plugs, i)
		}
	}

	return plugs
}

// Main implements the `buffalo info` command. Buffalo's checks
// are run first, then any plugins that implement plugins.Informer
// will be run in order at the end.
func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	return nil
}
