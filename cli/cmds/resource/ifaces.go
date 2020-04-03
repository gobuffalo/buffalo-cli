package resource

import (
	"context"
	"flag"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/spf13/pflag"
)

type BeforeGenerator interface {
	plugins.Plugin
	BeforeGenerateResource(ctx context.Context, root string, args []string) error
}

type AfterGenerator interface {
	plugins.Plugin
	AfterGenerateResource(ctx context.Context, root string, args []string, err error) error
}

type Actioner interface {
	plugins.Plugin
	GenerateResourceActions(ctx context.Context, root string, args []string) error
}

type ActionTester interface {
	plugins.Plugin
	GenerateResourceActionTests(ctx context.Context, root string, args []string) error
}

type Modeler interface {
	plugins.Plugin
	GenerateResourceModels(ctx context.Context, root string, args []string) error
}

type ModelTester interface {
	plugins.Plugin
	GenerateResourceModelTests(ctx context.Context, root string, args []string) error
}

type Templater interface {
	plugins.Plugin
	GenerateResourceTemplates(ctx context.Context, root string, args []string) error
}

type TemplateTester interface {
	plugins.Plugin
	GenerateResourceTemplateTests(ctx context.Context, root string, args []string) error
}

type Migrationer interface {
	plugins.Plugin
	GenerateResourceMigrations(ctx context.Context, root string, args []string) error
}

type MigrationTester interface {
	plugins.Plugin
	GenerateResourceMigrationTests(ctx context.Context, root string, args []string) error
}

type Flagger interface {
	plugins.Plugin
	ResourceFlags() []*flag.Flag
}

type Pflagger interface {
	plugins.Plugin
	ResourceFlags() []*pflag.Flag
}

type Stdouter = plugio.Outer
