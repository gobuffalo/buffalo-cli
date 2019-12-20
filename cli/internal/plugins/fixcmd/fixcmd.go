package fixcmd

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd/fix"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugins.Plugin = &FixCmd{}
var _ plugins.PluginNeeder = &FixCmd{}
var _ plugins.PluginScoper = &FixCmd{}
var _ plugprint.Describer = &FixCmd{}
var _ plugprint.SubCommander = &FixCmd{}

type FixCmd struct {
	pluginsFn plugins.PluginFeeder
}

func (fc *FixCmd) WithPlugins(f plugins.PluginFeeder) {
	fc.pluginsFn = f
}

func (fc *FixCmd) Name() string {
	return "fix"
}

func (fc *FixCmd) Description() string {
	return "Attempt to fix a Buffalo application's API to match version in go.mod"
}

func (f FixCmd) String() string {
	return f.Name()
}

// Fix runs any Fixers that are in the Plugins.
// If no arguments are provided it will run all fixers in the Plugins.
// Otherwise Fix will run the fixers for the arguments provided.
// 	buffalo fix
// 	buffalo fix plush pop
// 	buffalo fix -h
func (fc *FixCmd) plugins(ctx context.Context, args []string) error {
	plugs := fc.ScopedPlugins()
	if len(args) > 0 {
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
			if err := f.Fix(ctx, []string{}); err != nil {
				return err
			}
		}
		return nil
	}

	for _, p := range plugs {
		f, ok := p.(Fixer)
		if !ok {
			continue
		}

		if err := f.Fix(ctx, args); err != nil {
			return err
		}
	}
	return nil
}

func (fc *FixCmd) Main(ctx context.Context, args []string) error {
	var help bool
	flags := pflag.NewFlagSet(fc.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVarP(&fix.YesToAll, "yes", "y", false, "update all without asking for confirmation")
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
		return fc.plugins(ctx, args)
	}

	if err := fix.Run(); err != nil {
		return err
	}
	return fc.plugins(ctx, args)
}

func (fc *FixCmd) SubCommands() []plugins.Plugin {
	return fc.ScopedPlugins()
}

func (fc *FixCmd) ScopedPlugins() []plugins.Plugin {
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
