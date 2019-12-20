package pop

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/pop/v5/soda/cmd"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/parser"
)

const filePath = "/database.yml"

type Buffalo struct{}

func (b *Buffalo) BuildVersion(ctx context.Context, root string) (string, error) {
	return cmd.Version, nil
}

func (p *Buffalo) PkgerDecls() (parser.Decls, error) {
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
	return "pop"
}

func (p *Buffalo) BuiltInit(ctx context.Context, args []string) error {
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

var _ buildcmd.BuildImporter = Buffalo{}

func (Buffalo) BuildImports(ctx context.Context, root string) ([]string, error) {
	dir := filepath.Join(root, "models")
	info, err := here.Dir(dir)
	if err != nil {
		return nil, nil
	}
	return []string{info.ImportPath}, nil
}
