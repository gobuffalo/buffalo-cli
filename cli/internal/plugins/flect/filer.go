package flect

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/plugins"
)

var _ build.PackFiler = &Filer{}
var _ plugins.Plugin = &Filer{}

type Filer struct{}

func (Filer) PluginName() string {
	return "flect/filer"
}

func (f *Filer) PackageFiles(ctx context.Context, root string) ([]string, error) {
	return []string{
		filepath.Join(root, filePath),
	}, nil
}
