package flect

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/flect"
	"github.com/markbates/pkger"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Filer{},
	}
}

const filePath = "/inflections.json"

type Filer struct{}

var _ build.PackFiler = &Filer{}

func (f *Filer) PackageFiles(ctx context.Context, root string) ([]string, error) {
	return []string{
		filepath.Join(root, filePath),
	}, nil
}

var _ plugins.Plugin = &Filer{}

func (Filer) Name() string {
	return "flect"
}

func (fl *Filer) BuiltInit(ctx context.Context, args []string) error {
	f, err := pkger.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to load inflections %s", err)
	}
	defer f.Close()

	err = flect.LoadInflections(f)
	if err != nil {
		return err
	}
	return nil
}
