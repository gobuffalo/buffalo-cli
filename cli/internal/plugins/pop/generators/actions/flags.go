package actions

import (
	"io"

	"github.com/spf13/pflag"
)

func (g *Generator) PrintFlags(w io.Writer) error {
	flags := g.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (g *Generator) ResourceFlags() []*pflag.Flag {
	var values []*pflag.Flag
	flags := g.Flags()
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	return values
}

func (g *Generator) Flags() *pflag.FlagSet {
	if g.flags != nil && g.flags.Parsed() {
		return g.flags
	}

	flags := pflag.NewFlagSet(g.PluginName(), pflag.ContinueOnError)

	flags.StringVarP(&g.ModelName, "model-name", "n", "", "name of the model to use [defaults to resource name]")
	flags.StringVarP(&g.ModelsPkg, "model-pkg", "p", "", "full import path of models package [default is <module>/models]")
	flags.StringVarP(&g.ModelsPkgSel, "model-pkg-sel", "s", "", "selector for the models package [default is path.Base(model-pkg)]")

	g.flags = flags
	return g.flags
}
