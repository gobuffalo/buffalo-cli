package buildtest

import "context"

type Builder func(ctx context.Context, root string, args []string) error

func (Builder) PluginName() string {
	return "buildtest/builder"
}

func (b Builder) Build(ctx context.Context, root string, args []string) error {
	return b(ctx, root, args)
}
