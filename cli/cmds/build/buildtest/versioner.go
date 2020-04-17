package buildtest

import "context"

type Versioner func(ctx context.Context, root string) (string, error)

func (Versioner) PluginName() string {
	return "buildtest/versioner"
}

func (v Versioner) BuildVersion(ctx context.Context, root string) (string, error) {
	return v(ctx, root)
}
