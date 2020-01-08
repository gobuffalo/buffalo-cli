package resourcegen

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/spf13/pflag"
)

type Generator struct {
}

func (g *Generator) Name() string {
	return "plush/resource"
}

func (g *Generator) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(g.Name(), pflag.ContinueOnError)

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

var _ resource.Templater = &Generator{}

func (g *Generator) GenerateResourceTemplates(ctx context.Context, root string, args []string) error {
	flags := g.Flags()

	var model string
	flags.StringVarP(&model, "use-model", "", "", "tells resource generator to reference an existing model in generated code")

	if err := flags.Parse(args); err != nil {
		return err
	}

	// args = flags.Args()

	return nil
}
