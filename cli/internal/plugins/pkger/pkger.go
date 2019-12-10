package pkger

import (
	"context"

	"github.com/markbates/pkger/cmd/pkger/cmds"
	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/parser"
)

type Buffalo struct{}

func (b *Buffalo) Package(ctx context.Context, root string) error {
	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	decls, err := parser.Parse(info)
	if err != nil {
		return err
	}

	if err := cmds.Package(info, "pkged.go", decls); err != nil {
		return err
	}

	return nil
}

func (b Buffalo) Name() string {
	return "pkger"
}
