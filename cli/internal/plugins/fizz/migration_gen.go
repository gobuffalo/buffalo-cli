package fizz

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/generate"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/pop/v5/genny/fizz/ctable"
	"github.com/gobuffalo/pop/v5/soda/cmd"
	"github.com/spf13/pflag"
)

var _ generate.Generator = &MigrationGen{}
var _ plugins.Plugin = MigrationGen{}
var _ plugprint.Describer = MigrationGen{}
var _ plugprint.FlagPrinter = &MigrationGen{}
var _ plugprint.NamedCommand = MigrationGen{}
var _ resource.Migrationer = &MigrationGen{}
var _ resource.Pflagger = &MigrationGen{}

type MigrationGen struct {
	Env           string
	MigrationType string
	Path          string
	TableName     string

	flags *pflag.FlagSet
}

func (MigrationGen) Name() string {
	return "fizz/migration"
}

func (MigrationGen) CmdName() string {
	return "migration"
}

func (MigrationGen) Description() string {
	return "Generate a fizz migration"
}

func (mg *MigrationGen) Generate(ctx context.Context, args []string) error {
	args = append([]string{"generate", "migration"}, args...)
	cmd.RootCmd.SetArgs(args)
	return cmd.RootCmd.Execute()
}

func (mg *MigrationGen) GenerateResourceMigrations(ctx context.Context, root string, args []string) error {
	path := mg.Path
	if len(path) == 0 {
		path = filepath.Join(root, "migrations")
	}

	env := mg.Env
	if len(env) == 0 {
		env = "development"
	}

	mt := mg.MigrationType
	if len(mt) == 0 {
		mt = "fizz"
	}

	var translator fizz.Translator
	if mt == "sql" {
		db, err := pop.Connect(env)
		if err != nil {
			return err
		}
		translator = db.Dialect.FizzTranslator()
	}

	atts, err := attrs.ParseArgs(args[1:]...)
	if err != nil {
		return err
	}

	nm := mg.TableName
	if len(nm) == 0 {
		nm = args[0]
	}
	model := name.New(nm)

	g, err := ctable.New(&ctable.Options{
		TableName:              model.Tableize().String(),
		Attrs:                  atts,
		Path:                   path,
		Type:                   mt,
		Translator:             translator,
		ForceDefaultID:         true,
		ForceDefaultTimestamps: true,
	})
	if err != nil {
		return err
	}

	run := genny.WetRunner(ctx)
	run.With(g)
	return run.Run()
}
