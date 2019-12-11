package pkger

import (
	"context"
	"os"

	"github.com/gobuffalo/here/there"
	"github.com/markbates/pkger/cmd/pkger/cmds"
	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/parser"
)

type Buffalo struct{}

func (b *Buffalo) Package(ctx context.Context, root string) error {
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

	os.RemoveAll("pkged.go")
	if err := cmds.Package(info, "pkged.go", decls); err != nil {
		return err
	}

	return nil
}

func (b Buffalo) Name() string {
	return "pkger"
}
