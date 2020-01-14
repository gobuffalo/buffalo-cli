package pop

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/pop/v5/soda/cmd"
)

const filePath = "/database.yml"

type Builder struct{}

var _ plugins.Plugin = Builder{}

func (Builder) Name() string {
	return "pop/builder"
}

var _ buildcmd.Versioner = &Builder{}

func (b *Builder) BuildVersion(ctx context.Context, root string) (string, error) {
	return cmd.Version, nil
}

var _ buildcmd.PackFiler = &Builder{}

func (b *Builder) PackageFiles(ctx context.Context, root string) ([]string, error) {
	return []string{
		filepath.Join(root, filePath),
	}, nil
}

var _ buildcmd.Importer = Builder{}

func (Builder) BuildImports(ctx context.Context, root string) ([]string, error) {
	dir := filepath.Join(root, "models")
	info, err := here.Dir(dir)
	if err != nil {
		return nil, nil
	}
	return []string{info.ImportPath}, nil
}
