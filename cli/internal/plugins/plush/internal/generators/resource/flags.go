package resource

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/spf13/pflag"
)

func (g *Generator) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(g.Name(), pflag.ContinueOnError)
	flags.StringVarP(&g.modelName, "model-name", "n", "", "name of the model to use [defaults to resource name]")
	return flags
}

var _ resource.Pflagger = &Generator{}

func (g *Generator) ResourceFlags() []*pflag.Flag {
	var values []*pflag.Flag
	flags := g.Flags()
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	return values
}

var _ plugprint.FlagPrinter = &Generator{}

func (g *Generator) PrintFlags(w io.Writer) error {
	flags := g.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
