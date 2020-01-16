package fizz

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/generatecmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/soda"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/pop/v5/genny/fizz/ctable"
)

type MigrationGen struct {
	path          string
	env           string
	migrationType string
	tableName     string
}

var _ plugins.Plugin = MigrationGen{}

func (MigrationGen) Name() string {
	return "fizz/migration"
}

var _ plugprint.NamedCommand = MigrationGen{}

func (MigrationGen) CmdName() string {
	return "migration"
}

var _ plugprint.Describer = MigrationGen{}

func (MigrationGen) Description() string {
	return "Generate a fizz migration"
}

var _ generatecmd.Generator = &MigrationGen{}

func (mg *MigrationGen) Generate(ctx context.Context, args []string) error {
	args = append([]string{"generate", "migration"}, args...)
	return soda.Main(ctx, args)
}

var _ resource.Migrationer = &MigrationGen{}

func (mg *MigrationGen) GenerateResourceMigrations(ctx context.Context, root string, args []string) error {
	path := mg.path
	if len(path) == 0 {
		path = filepath.Join(root, "migrations")
	}

	env := mg.env
	if len(env) == 0 {
		env = "development"
	}

	mt := mg.migrationType
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

	nm := mg.tableName
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
