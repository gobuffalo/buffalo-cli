package info

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/spf13/pflag"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Cmd{},
	}
}

var _ plugins.Plugin = &Cmd{}
var _ plugins.PluginNeeder = &Cmd{}
var _ plugins.PluginScoper = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}

type Cmd struct {
	pluginsFn plugins.PluginFeeder
	help      bool
}

func (ic *Cmd) WithPlugins(f plugins.PluginFeeder) {
	ic.pluginsFn = f
}

func (ic *Cmd) PrintFlags(w io.Writer) error {
	flags := ic.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (ic *Cmd) Name() string {
	return "info"
}

func (ic *Cmd) Description() string {
	return "Print diagnostic information (useful for debugging)"
}

func (i Cmd) String() string {
	return i.Name()
}

// Info runs all of the plugins that implement the
// `Informer` interface in order.
func (ic *Cmd) plugins(ctx context.Context, args []string) error {
	for _, p := range ic.ScopedPlugins() {
		i, ok := p.(Informer)
		if !ok {
			continue
		}
		if err := i.Info(ctx, args); err != nil {
			return err
		}
	}
	return nil
}

func (ic *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin

	if ic.pluginsFn == nil {
		return plugs
	}
	for _, p := range ic.pluginsFn() {
		if i, ok := p.(Informer); ok {
			plugs = append(plugs, i)
		}
	}

	return plugs
}

func (ic *Cmd) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(ic.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVarP(&ic.help, "help", "h", false, "print this help")
	return flags
}

// Main implements the `buffalo info` command. Buffalo's checks
// are run first, then any plugins that implement plugins.Informer
// will be run in order at the end.
func (ic *Cmd) Main(ctx context.Context, args []string) error {
	return nil
}
