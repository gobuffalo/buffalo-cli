package fizz

import (
	"io"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugprint.FlagPrinter = &MigrationGen{}

func (g *MigrationGen) PrintFlags(w io.Writer) error {
	flags := g.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

var _ resource.Pflagger = &MigrationGen{}

func (g *MigrationGen) ResourceFlags() []*pflag.Flag {
	var values []*pflag.Flag
	flags := g.Flags()
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	return values
}

func (g *MigrationGen) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(g.Name(), pflag.ContinueOnError)

	flags.StringVarP(&g.env, "env", "e", "development", "The environment you want to run migrations against. Will use $GO_ENV if set.")
	flags.StringVarP(&g.migrationType, "type", "", "fizz", "sets the type of migration files for model (sql or fizz)")
	flags.StringVarP(&g.path, "path", "p", "./migrations", "Path to the migrations folder")
	flags.StringVarP(&g.tableName, "table-name", "", "", "name for the database table [defaults to pluralized model name]")
	return flags
}
