package resource

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/v2/internal/flagger"
	"github.com/gobuffalo/plugins/plugflag"
	"github.com/spf13/pflag"
)

func (g *Generator) PrintFlags(w io.Writer) error {
	flags := g.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (g *Generator) Flags() *pflag.FlagSet {
	if g.flags != nil && g.flags.Parsed() {
		return g.flags
	}

	flags := pflag.NewFlagSet(g.PluginName(), pflag.ContinueOnError)
	flags.BoolVar(&g.SkipActionTests, "skip-action-tests", false, "skip generating action tests")
	flags.BoolVar(&g.SkipActions, "skip-actions", false, "skip generating actions")
	flags.BoolVar(&g.SkipMigrationTests, "skip-migration-tests", false, "skip generating migration tests")
	flags.BoolVar(&g.SkipMigrations, "skip-migrations", false, "skip generating migrations")
	flags.BoolVar(&g.SkipModelTests, "skip-model-tests", false, "skip generating model tests")
	flags.BoolVar(&g.SkipModels, "skip-models", false, "skip generating models")
	flags.BoolVar(&g.SkipTemplateTests, "skip-template-tests", false, "skip generating template tests")
	flags.BoolVar(&g.SkipTemplates, "skip-templates", false, "skip generating templates")
	flags.BoolVarP(&g.help, "help", "h", false, "print this help")

	plugs := g.ScopedPlugins()

	for _, p := range plugs {
		switch t := p.(type) {
		case Flagger:
			for _, f := range plugflag.Clean(p, t.ResourceFlags()) {
				flags.AddGoFlag(f)
			}
		case Pflagger:
			for _, f := range flagger.CleanPflags(p, t.ResourceFlags()) {
				flags.AddGoFlag(f)
			}
		}
	}
	g.flags = flags
	return g.flags
}
