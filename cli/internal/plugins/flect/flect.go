package flect

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/flect"
	"github.com/markbates/pkger"
)

const filePath = "/inflections.json"

var _ plugins.Plugin = &Buffalo{}
var _ buildcmd.PackFiler = &Buffalo{}

type Buffalo struct{}

func (f *Buffalo) PackageFiles(ctx context.Context, root string) ([]string, error) {
	return []string{
		filepath.Join(root, filePath),
	}, nil
}

func (Buffalo) Name() string {
	return "flect"
}

func (fl *Buffalo) BuiltInit(ctx context.Context, args []string) error {
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
