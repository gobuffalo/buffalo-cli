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
	flags := pflag.NewFlagSet(g.Name(), pflag.ContinueOnError)

	flags.StringVarP(&g.path, "path", "p", "./migrations", "Path to the migrations folder")
	flags.StringVarP(&g.env, "env", "e", "development", "The environment you want to run migrations against. Will use $GO_ENV if set.")

	flags.StringVarP(&g.migrationType, "migration-type", "", "fizz", "sets the type of migration files for model (sql or fizz)")

	flags.StringVarP(&g.tableName, "table-name", "", "", "name for the database table [defaults to pluralized model name]")
	// flags.BoolVarP(&g.help, "help", "h", false, "print this help")

	// plugs := g.ScopedPlugins()
	//
	// for _, p := range plugs {
	// 	switch t := p.(type) {
	// 	case Flagger:
	// 		for _, f := range t.ResourceFlags() {
	// 			flags.AddGoFlag(f)
	// 		}
	// 	case Pflagger:
	// 		for _, f := range t.ResourceFlags() {
	// 			flags.AddFlag(f)
	// 		}
	// 	}
	// }
	return flags
}
