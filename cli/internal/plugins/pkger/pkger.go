package pkger

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here/there"
	"github.com/markbates/pkger/cmd/pkger/cmds"
	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/parser"
)

var _ plugins.Plugin = &Buffalo{}
var _ plugprint.WithPlugins = &Buffalo{}

type Buffalo struct {
	Plugins func() plugins.Plugins
}

type Decler interface {
	PkgerDecls() (parser.Decls, error)
}

func (b *Buffalo) WithPlugins() plugins.Plugins {
	var plugs plugins.Plugins
	if b.Plugins != nil {
		plugs = b.Plugins()
	}

	var builders plugins.Plugins
	for _, p := range plugs {
		switch p.(type) {
		case Decler:
			builders = append(builders, p)
		}
	}
	return builders
}

func (b *Buffalo) Build(ctx context.Context, args []string) error {
	return b.Package(ctx, ".")
}

func (b *Buffalo) Package(ctx context.Context, root string) error {
	defer func() {
		if err := recover(); err != nil {
			os.RemoveAll("pkged.go")
		}
	}()
	thar, err := there.Current()
	if err != nil {
		return err
	}
	info := here.Info{
		Dir:        thar.Dir,
		ImportPath: thar.ImportPath,
		Name:       thar.Name,
		Module: here.Module{
			Path:      thar.Module.Path,
			Main:      thar.Module.Main,
			Dir:       thar.Module.Dir,
			GoMod:     thar.Module.GoMod,
			GoVersion: thar.Module.GoVersion,
		},
	}

	decls, err := parser.Parse(info)
	if err != nil {
		return err
	}

	for _, p := range b.WithPlugins() {
		pd, ok := p.(Decler)
		if !ok {
			continue
		}
		ds, err := pd.PkgerDecls()
		if err != nil {
			return err
		}
		decls = append(decls, ds...)
	}

	os.RemoveAll("pkged.go")
	if err := cmds.Package(info, "pkged.go", decls); err != nil {
		return err
	}

	return nil
}

func (b Buffalo) Name() string {
	return "pkger"
}
