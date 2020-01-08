package resource

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

func (g *Generator) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(g.Name(), pflag.ContinueOnError)
	flags.BoolVarP(&g.help, "help", "h", false, "print this help")
	// flags.StringVarP(&g.model, "use-model", "", "", "tells resource generator to reference an existing model in generated code")

	plugs := g.ScopedPlugins()

	for _, p := range plugs {
		switch t := p.(type) {
		case Flagger:
			for _, f := range t.ResourceFlags() {
				flags.AddGoFlag(f)
			}
		case Pflagger:
			for _, f := range t.ResourceFlags() {
				flags.AddFlag(f)
			}
		}
	}
	return flags
}
