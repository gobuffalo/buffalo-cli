package fizz

import (
	"io"

	"github.com/spf13/pflag"
)

func (g *MigrationGen) PrintFlags(w io.Writer) error {
	flags := g.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (g *MigrationGen) ResourceFlags() []*pflag.Flag {
	var values []*pflag.Flag
	flags := g.Flags()
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	return values
}

func (g *MigrationGen) Flags() *pflag.FlagSet {
	if g.flags != nil {
		return g.flags
	}

	flags := pflag.NewFlagSet(g.PluginName(), pflag.ContinueOnError)

	flags.StringVarP(&g.Env, "env", "e", "development", "The environment you want to run migrations against. Will use $GO_ENV if set.")
	flags.StringVarP(&g.MigrationType, "type", "", "fizz", "sets the type of migration files for model (sql or fizz)")
	flags.StringVarP(&g.Path, "path", "p", "./migrations", "Path to the migrations folder")
	flags.StringVarP(&g.TableName, "table-name", "", "", "name for the database table [defaults to pluralized model name]")

	g.flags = flags
	return g.flags
}
