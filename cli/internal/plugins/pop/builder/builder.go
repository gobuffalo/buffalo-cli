package builder

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/pop/v5/soda/cmd"
)

var _ build.Importer = Builder{}
var _ build.PackFiler = &Builder{}
var _ build.Versioner = &Builder{}
var _ plugins.Plugin = Builder{}

const filePath = "/database.yml"

type Builder struct{}

func (Builder) Name() string {
	return "pop/builder"
}

func (b *Builder) BuildVersion(ctx context.Context, root string) (string, error) {
	return cmd.Version, nil
}

func (b *Builder) PackageFiles(ctx context.Context, root string) ([]string, error) {
	return []string{
		filepath.Join(root, filePath),
	}, nil
}

func (Builder) BuildImports(ctx context.Context, root string) ([]string, error) {
	dir := filepath.Join(root, "models")
	info, err := here.Dir(dir)
	if err != nil {
		return nil, nil
	}
	return []string{info.ImportPath}, nil
}
