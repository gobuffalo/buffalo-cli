package generate

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/v2/internal/flagger"
	"github.com/spf13/pflag"
)

func (g *Migration) PrintFlags(w io.Writer) error {
	flags := g.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (g *Migration) ResourceFlags() []*pflag.Flag {
	return flagger.SetToSlice(g.Flags())
}

func (g *Migration) Flags() *pflag.FlagSet {
	if g.flags != nil {
		return g.flags
	}

	flags := pflag.NewFlagSet(g.PluginName(), pflag.ContinueOnError)

	flags.StringVarP(&g.env, "env", "e", "development", "The environment you want to run migrations against. Will use $GO_ENV if set.")
	flags.StringVarP(&g.migrationType, "type", "", "fizz", "sets the type of migration files for model (sql or fizz)")
	flags.StringVarP(&g.path, "path", "p", "./migrations", "Path to the migrations folder")
	flags.StringVarP(&g.tableName, "table-name", "", "", "name for the database table [defaults to pluralized model name]")

	g.flags = flags
	return g.flags
}
