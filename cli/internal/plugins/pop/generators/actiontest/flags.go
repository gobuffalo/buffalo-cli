package actiontest

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
	if g.flags != nil {
		return g.flags
	}

	flags := pflag.NewFlagSet(g.PluginName(), pflag.ContinueOnError)

	flags.StringVarP(&g.TestPkg, "test-pkg", "t", "", "name of the test package to use [default 'actions']")

	g.flags = flags
	return g.flags
}
