package buildtest

import "context"

type AfterBuilder func(ctx context.Context, root string, args []string, err error) error

func (AfterBuilder) PluginName() string {
	return "buildtest/builder"
}

func (b AfterBuilder) AfterBuild(ctx context.Context, root string, args []string, err error) error {
	return b(ctx, root, args, err)
}
