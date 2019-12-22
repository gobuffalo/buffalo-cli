package flect

import (
	"context"
	"fmt"

	pkgerplug "github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pkger"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/here"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/parser"
)

const filePath = "/inflections.json"

var _ plugins.Plugin = &Buffalo{}
var _ pkgerplug.Decler = &Buffalo{}

type Buffalo struct{}

func (f *Buffalo) PkgerDecls() (parser.Decls, error) {
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

func (Buffalo) Name() string {
	return "flect"
}

func (fl *Buffalo) BuiltInit(ctx context.Context, args []string) error {
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
