package newapp

import (
	"context"
	"fmt"
	"io"

	"github.com/gobuffalo/buffalo-cli/v2/internal/flagger"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugflag"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

type Executer struct {
	pluginsFn plugins.Feeder
	flags     *pflag.FlagSet
	help      bool
}

func (Executer) PluginName() string {
	return "newapp/executer"
}

func (e *Executer) WithPlugins(f plugins.Feeder) {
	e.pluginsFn = f
}

func (e *Executer) ScopedPlugins() []plugins.Plugin {
	if e.pluginsFn == nil {
		return nil
	}

	var plugs []plugins.Plugin
	for _, p := range e.pluginsFn() {
		switch p.(type) {
		case Stdouter:
			plugs = append(plugs, p)
		case Stdiner:
			plugs = append(plugs, p)
		case Stderrer:
			plugs = append(plugs, p)
		case Newapper:
			plugs = append(plugs, p)
		case AfterNewapper:
			plugs = append(plugs, p)
		case Flagger:
			plugs = append(plugs, p)
		case Pflagger:
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (e *Executer) Execute(ctx context.Context, root string, name string, args []string) error {
	flags := e.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	plugs := e.ScopedPlugins()
	if e.help {
		return plugprint.Print(plugio.Stdout(plugs...), e)
	}

	var during []Newapper
	var after []AfterNewapper

	for _, p := range plugs {
		switch t := p.(type) {
		case Newapper:
			during = append(during, t)
		case AfterNewapper:
			after = append(after, t)
		}
	}

	var err error
	for _, p := range during {
		fmt.Println(">>>TODO DURING ", p.PluginName())
		if err = p.Newapp(ctx, root, name, args); err != nil {
			break
		}
	}

	for _, p := range after {
		fmt.Println(">>>TODO AFTER ", p.PluginName())
		if err := p.AfterNewapp(ctx, root, name, args, err); err != nil {
			return err
		}
	}
	return err
}

func (ex *Executer) PrintFlags(w io.Writer) error {
	flags := ex.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (ex *Executer) Flags() *pflag.FlagSet {
	if ex.flags != nil {
		return ex.flags
	}
	flags := (&Cmd{}).Flags()
	flags.BoolVarP(&ex.help, "help", "h", false, "print this help")

	plugs := ex.ScopedPlugins()

	for _, p := range plugs {
		switch t := p.(type) {
		case Flagger:
			for _, f := range plugflag.Clean(p, t.NewappFlags()) {
				flags.AddGoFlag(f)
			}
		case Pflagger:
			for _, f := range flagger.CleanPflags(p, t.NewappFlags()) {
				flags.AddFlag(f)
			}
		}
	}

	ex.flags = flags
	return ex.flags
}

func Execute(plugs []plugins.Plugin, ctx context.Context, root string, name string, args []string) error {
	e := &Executer{}
	e.WithPlugins(func() []plugins.Plugin {
		return plugs
	})
	return e.Execute(ctx, root, name, args)
}
