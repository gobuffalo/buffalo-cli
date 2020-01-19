package build

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here"
	"github.com/markbates/safe"
)

func (bc *Cmd) pack(ctx context.Context, info here.Info, plugs []plugins.Plugin) error {
	var files []string
	for _, p := range plugs {
		pkg, ok := p.(PackFiler)
		if !ok {
			continue
		}
		err := safe.RunE(func() error {
			res, err := pkg.PackageFiles(ctx, info.Dir)
			if err != nil {
				return err
			}
			files = append(files, res...)
			return nil
		})
		if err != nil {
			return err
		}
	}

	for _, p := range plugs {
		pkg, ok := p.(Packager)
		if !ok {
			continue
		}
		err := safe.RunE(func() error {
			return pkg.Package(ctx, info.Dir, files)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
