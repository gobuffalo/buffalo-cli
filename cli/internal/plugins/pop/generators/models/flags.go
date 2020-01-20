package models

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

	flags := pflag.NewFlagSet(g.Name(), pflag.ContinueOnError)

	flags.StringVarP(&g.ModelPath, "path", "", "models", "the path the model will be created in")
	flags.StringVarP(&g.ModelPkg, "pkg", "", "models", "the import part the model will be created in")
	flags.StringVarP(&g.StructTag, "struct-tag", "", "json", "sets the struct tags for model (xml/json/jsonapi)")

	g.flags = flags
	return g.flags
}
