package buildtest

import "context"

type Packager func(ctx context.Context, root string, args []string) error

func (Packager) PluginName() string {
	return "buildtest/packager"
}

func (p Packager) Package(ctx context.Context, root string, args []string) error {
	return p(ctx, root, args)
}
