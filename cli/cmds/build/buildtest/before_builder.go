package buildtest

import "context"

type BeforeBuilder func(ctx context.Context, root string, args []string) error

func (BeforeBuilder) PluginName() string {
	return "buildtest/builder"
}

func (b BeforeBuilder) BeforeBuild(ctx context.Context, root string, args []string) error {
	return b(ctx, root, args)
}
