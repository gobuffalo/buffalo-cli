package resource

import "context"

type BeforeGenerator interface {
	BeforeGenerateResource(ctx context.Context, root string, args []string) error
}

type ResourceGenerator interface {
	GenerateResource(ctx context.Context, root string, args []string) error
}

type AfterGenerator interface {
	AfterGenerateResource(ctx context.Context, root string, args []string, err error) error
}

type Actioner interface {
	GenerateResourceActions(ctx context.Context, root string, args []string) error
}

type ActionTester interface {
	GenerateResourceActionTests(ctx context.Context, root string, args []string) error
}

type Modeler interface {
	GenerateResourceModels(ctx context.Context, root string, args []string) error
}

type ModelTester interface {
	GenerateResourceModelTests(ctx context.Context, root string, args []string) error
}

type Templater interface {
	GenerateResourceTemplates(ctx context.Context, root string, args []string) error
}

type TemplateTester interface {
	GenerateResourceTemplateTests(ctx context.Context, root string, args []string) error
}

type Migrationer interface {
	GenerateResourceMigrations(ctx context.Context, root string, args []string) error
}

type MigrationTester interface {
	GenerateResourceMigrationTests(ctx context.Context, root string, args []string) error
}
