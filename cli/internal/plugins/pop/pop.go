package pop

import (
	"context"

	"github.com/gobuffalo/here/there"
	"github.com/gobuffalo/pop"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/parser"
)

const filePath = "/database.yml"

type Pop struct{}

func (p *Pop) PkgerDecls() (parser.Decls, error) {
	info, err := there.Current()
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

func (Pop) Name() string {
	return "pop"
}

func (p *Pop) BuiltInit(ctx context.Context, args []string) error {
	f, err := pkger.Open("/database.yml")
	if err != nil {
		return err
	}
	defer f.Close()

	err = pop.LoadFrom(f)
	if err != nil {
		return err
	}
	return nil
}
