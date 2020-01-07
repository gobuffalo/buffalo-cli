package fizz

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/generatecmd"
	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/resource"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	soda "github.com/gobuffalo/pop/v5/soda/cmd"
)

type MigrationGen struct{}

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
	soda.RootCmd.SetArgs(args)
	return soda.RootCmd.Execute()
}

var _ resource.Migrationer = &MigrationGen{}

func (mg *MigrationGen) GenerateResourceMigrations(ctx context.Context, root string, args []string) error {
	return mg.Generate(ctx, args)
}
