package flect

import (
	"context"
	"fmt"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/here"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/parser"
)

const filePath = "/inflections.json"

type Flect struct{}

func (f *Flect) PkgerDecls() (parser.Decls, error) {
	info, err := here.Current()
	if err != nil {
		return nil, err
	}

	var decls parser.Decls

	d, err := parser.NewInclude(info, filePath)
	if err != nil {
		return nil, err
	}
	decls = append(decls, d)

	return decls, nil
}

func (Flect) Name() string {
	return "flect"
}

func (fl *Flect) BuiltInit(ctx context.Context, args []string) error {
	f, err := pkger.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to load inflections %w", err)
	}
	defer f.Close()

	err = flect.LoadInflections(f)
	if err != nil {
		return err
	}
	return nil
}
