package resource

import (
	"io"

	"github.com/spf13/pflag"
)

func (g *Generator) Flags() *pflag.FlagSet {
	if g.flags != nil && g.flags.Parsed() {
		return g.flags
	}
	flags := pflag.NewFlagSet(g.PluginName(), pflag.ContinueOnError)
	flags.StringVarP(&g.modelName, "model-name", "n", "", "name of the model to use [defaults to resource name]")

	g.flags = flags
	return g.flags
}

func (g *Generator) ResourceFlags() []*pflag.Flag {
	var values []*pflag.Flag
	flags := g.Flags()
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	return values
}

func (g *Generator) PrintFlags(w io.Writer) error {
	flags := g.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
