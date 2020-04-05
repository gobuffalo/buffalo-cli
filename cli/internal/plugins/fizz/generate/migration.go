package generate

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/generate"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/resource"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/pop/v5/genny/fizz/ctable"
	"github.com/gobuffalo/pop/v5/soda/cmd"
	"github.com/spf13/pflag"
)

var _ generate.Generator = &Migration{}
var _ plugcmd.Namer = Migration{}
var _ plugins.Plugin = Migration{}
var _ plugprint.Describer = Migration{}
var _ plugprint.FlagPrinter = &Migration{}
var _ resource.Migrationer = &Migration{}
var _ resource.Pflagger = &Migration{}

type Migration struct {
	env           string
	migrationType string
	path          string
	tableName     string
	flags         *pflag.FlagSet
}

func (Migration) PluginName() string {
	return "fizz/migration"
}

func (Migration) CmdName() string {
	return "migration"
}

func (Migration) Description() string {
	return "Generate a fizz migration"
}

func (mg *Migration) Generate(ctx context.Context, root string, args []string) error {
	args = append([]string{"generate", "migration"}, args...)
	cmd.RootCmd.SetArgs(args)
	return cmd.RootCmd.Execute()
}

func (mg *Migration) GenerateResourceMigrations(ctx context.Context, root string, args []string) error {
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
