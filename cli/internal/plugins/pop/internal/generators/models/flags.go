package models

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugprint.FlagPrinter = &Generator{}

func (g *Generator) PrintFlags(w io.Writer) error {
	flags := g.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
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

func (g *Generator) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(g.Name(), pflag.ContinueOnError)

	flags.StringVarP(&g.modelPath, "path", "", "models", "the path the model will be created in")
	flags.StringVarP(&g.modelPkg, "pkg", "", "models", "the import part the model will be created in")
	flags.StringVarP(&g.structTag, "struct-tag", "", "json", "sets the struct tags for model (xml/json/jsonapi)")
	return flags
}
