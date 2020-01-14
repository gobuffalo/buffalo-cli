package grifts

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
	grifts "github.com/markbates/grift/cmd"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Buffalo{},
		&Generator{},
	}
}

type Buffalo struct {
}

var _ plugprint.Aliases = Buffalo{}

func (Buffalo) Aliases() []string {
	return []string{"task", "tasks", "t"}
}

func (bc *Buffalo) Main(ctx context.Context, args []string) error {
	return grifts.Run("buffalo grifts", args)
}

var _ plugins.Plugin = Buffalo{}

func (Buffalo) Name() string {
	return "grifts"
}

var _ buildcmd.Importer = Buffalo{}

func (Buffalo) BuildImports(ctx context.Context, root string) ([]string, error) {
	dir := filepath.Join(root, "grifts")
	info, err := here.Dir(dir)
	if err != nil {
		return nil, nil
	}
	return []string{info.ImportPath}, nil
}
