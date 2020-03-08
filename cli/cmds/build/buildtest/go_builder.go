package buildtest

import "context"

type GoBuilder func(ctx context.Context, root string, args []string) error

func (GoBuilder) PluginName() string {
	return "buildtest/go-builder"
}
func (g GoBuilder) GoBuild(ctx context.Context, root string, args []string) error {
	return g(ctx, root, args)
}
