package grifts

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/here"
)

type Buffalo struct {
}

func (Buffalo) Aliases() []string {
	return []string{"task", "t"}
}

func (Buffalo) Name() string {
	return "grifts"
}

var _ buildcmd.BuildImporter = Buffalo{}

func (Buffalo) BuildImports(ctx context.Context, root string) ([]string, error) {
	dir := filepath.Join(root, "grifts")
	info, err := here.Dir(dir)
	if err != nil {
		return nil, nil
	}
	return []string{info.ImportPath}, nil
}
