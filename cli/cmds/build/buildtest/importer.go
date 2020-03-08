package buildtest

import "context"

type Importer func(ctx context.Context, root string) ([]string, error)

func (Importer) PluginName() string {
	return "buildtest/importer"
}

func (i Importer) BuildImports(ctx context.Context, root string) ([]string, error) {
	return i(ctx, root)
}
