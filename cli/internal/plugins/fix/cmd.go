package fix

import (
	"context"
	"fmt"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/markbates/safe"
	"github.com/spf13/pflag"
)

var _ plugcmd.SubCommander = &Cmd{}
var _ plugins.Plugin = &Cmd{}
var _ plugins.Needer = &Cmd{}
var _ plugins.Scoper = &Cmd{}
var _ plugprint.Describer = &Cmd{}

type Cmd struct {
	flags     *pflag.FlagSet
	help      bool
	pluginsFn plugins.Feeder
}

func (fc *Cmd) WithPlugins(f plugins.Feeder) {
	fc.pluginsFn = f
}

func (fc *Cmd) PluginName() string {
	return "fix"
}

func (fc *Cmd) Description() string {
	return "Attempt to fix a Buffalo application's API to match version in go.mod"
}

func (f Cmd) String() string {
	return f.PluginName()
}

// Fix runs any Fixers that are in the Plugins.
// If no arguments are provided it will run all fixers in the Plugins.
// Otherwise Fix will run the fixers for the arguments provided.
// 	buffalo fix
// 	buffalo fix plush pop
// 	buffalo fix -h
func (fc *Cmd) fixPlugins(ctx context.Context, root string, args []string) error {
	plugs := fc.ScopedPlugins()

	if len(args) == 0 {
		for _, p := range plugs {
			f, ok := p.(Fixer)
			if !ok {
				continue
			}
			err := safe.RunE(func() error {
				return f.Fix(ctx, root, args)
			})
			if err != nil {
				return err
			}
		}
		return nil
	}

	fixers := map[string]Fixer{}
	for _, p := range plugs {
		f, ok := p.(Fixer)
		if !ok {
			continue
		}

		fixers[p.PluginName()] = f
	}

	for _, a := range args {
		f, ok := fixers[a]
		if !ok {
			return fmt.Errorf("unknown fixer %s", a)
		}
		err := safe.RunE(func() error {
			return f.Fix(ctx, root, []string{})
		})
		if err != nil {
			return err
		}
	}
	return nil

}

func (fc *Cmd) Main(ctx context.Context, root string, args []string) error {
	flags := fc.Flags()

	if err := flags.Parse(args); err != nil {
		return err
	}

	if fc.help {
		return plugprint.Print(plugio.Stdout(fc.ScopedPlugins()...), fc)
	}

	if len(args) > 0 {
		return fc.fixPlugins(ctx, root, args)
	}

	return fc.fixPlugins(ctx, root, args)
}

func (fc *Cmd) SubCommands() []plugins.Plugin {
	return fc.ScopedPlugins()
}

func (fc *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if fc.pluginsFn == nil {
		return plugs
	}

	for _, p := range fc.pluginsFn() {
		switch p.(type) {
		case Fixer:
			plugs = append(plugs, p)
		case Flagger:
			plugs = append(plugs, p)
		case Pflagger:
			plugs = append(plugs, p)
		case Stdouter:
			plugs = append(plugs, p)
		}
	}
	return plugs
}
