package buildtest

import "context"

type PackFiler func(ctx context.Context, root string) ([]string, error)

func (PackFiler) PluginName() string {
	return "buildtest/packfiler"
}

func (p PackFiler) PackageFiles(ctx context.Context, root string) ([]string, error) {
	return p(ctx, root)
}
