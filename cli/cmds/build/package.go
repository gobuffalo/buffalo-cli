package build

import (
	"context"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
)

func (bc *Cmd) pack(ctx context.Context, info here.Info, plugs []plugins.Plugin) error {
	var files []string
	for _, p := range plugs {
		pkg, ok := p.(PackFiler)
		if !ok {
			continue
		}
		res, err := pkg.PackageFiles(ctx, info.Dir)
		if err != nil {
			return plugins.Wrap(p, err)
		}
		files = append(files, res...)
	}

	for _, p := range plugs {
		if pkg, ok := p.(Packager); ok {

			if err := pkg.Package(ctx, info.Dir, files); err != nil {
				return plugins.Wrap(p, err)
			}
		}
	}
	return nil
}
