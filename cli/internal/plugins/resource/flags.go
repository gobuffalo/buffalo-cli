package resource

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/internal/flagger"
	"github.com/gobuffalo/buffalo-cli/plugins"
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
	flags.BoolVar(&g.skipActionTests, "skip-action-tests", false, "skip generating action tests")
	flags.BoolVar(&g.skipActions, "skip-actions", false, "skip generating actions")
	flags.BoolVar(&g.skipMigrationTests, "skip-migration-tests", false, "skip generating migration tests")
	flags.BoolVar(&g.skipMigrations, "skip-migrations", false, "skip generating migrations")
	flags.BoolVar(&g.skipModelTests, "skip-model-tests", false, "skip generating model tests")
	flags.BoolVar(&g.skipModels, "skip-models", false, "skip generating models")
	flags.BoolVar(&g.skipTemplateTests, "skip-template-tests", false, "skip generating template tests")
	flags.BoolVar(&g.skipTemplates, "skip-templates", false, "skip generating templates")
	flags.BoolVarP(&g.help, "help", "h", false, "print this help")

	plugs := g.ScopedPlugins()

	for _, p := range plugs {
		switch t := p.(type) {
		case Flagger:
			for _, f := range plugins.CleanFlags(p, t.ResourceFlags()) {
				flags.AddGoFlag(f)
			}
		case Pflagger:
			for _, f := range flagger.CleanPflags(p, t.ResourceFlags()) {
				flags.AddGoFlag(f)
			}
		}
	}
	return flags
}
