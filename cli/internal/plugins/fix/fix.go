package fix

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/markbates/safe"
	"github.com/spf13/pflag"
)

type Cmd struct {
	pluginsFn plugins.PluginFeeder
}

var _ plugins.PluginNeeder = &Cmd{}

func (fc *Cmd) WithPlugins(f plugins.PluginFeeder) {
	fc.pluginsFn = f
}

var _ plugins.Plugin = &Cmd{}

func (fc *Cmd) Name() string {
	return "fix"
}

var _ plugprint.Describer = &Cmd{}

func (fc *Cmd) Description() string {
	return "Attempt to fix a Buffalo application's API to match version in go.mod"
}

func (f Cmd) String() string {
	return f.Name()
}

// Fix runs any Fixers that are in the Plugins.
// If no arguments are provided it will run all fixers in the Plugins.
// Otherwise Fix will run the fixers for the arguments provided.
// 	buffalo fix
// 	buffalo fix plush pop
// 	buffalo fix -h
func (fc *Cmd) fixPlugins(ctx context.Context, args []string) error {
	plugs := fc.ScopedPlugins()

	if len(args) == 0 {
		for _, p := range plugs {
			f, ok := p.(Fixer)
			if !ok {
				continue
			}
			err := safe.RunE(func() error {
				return f.Fix(ctx, args)
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

		fixers[p.Name()] = f
	}

	for _, a := range args {
		f, ok := fixers[a]
		if !ok {
			return fmt.Errorf("unknown fixer %s", a)
		}
		err := safe.RunE(func() error {
			return f.Fix(ctx, []string{})
		})
		if err != nil {
			return err
		}
	}
	return nil

}

func (fc *Cmd) Main(ctx context.Context, args []string) error {
	var help bool
	flags := pflag.NewFlagSet(fc.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVarP(&help, "help", "h", false, "print this help")

	if err := flags.Parse(args); err != nil {
		return err
	}

	ioe := plugins.CtxIO(ctx)
	out := ioe.Stdout()

	if help {
		return plugprint.Print(out, fc)
	}

	if len(args) > 0 {
		return fc.fixPlugins(ctx, args)
	}

	return fc.fixPlugins(ctx, args)
}

var _ plugprint.SubCommander = &Cmd{}

func (fc *Cmd) SubCommands() []plugins.Plugin {
	return fc.ScopedPlugins()
}

var _ plugins.PluginScoper = &Cmd{}

func (fc *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if fc.pluginsFn != nil {
		for _, p := range fc.pluginsFn() {
			if _, ok := p.(Fixer); ok {
				plugs = append(plugs, p)
			}
		}
	}
	return plugs
}
