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

func (Buffalo) Name() string {
	return "pop"
}

// Main adds the `buffalo pop` sub-command.
func (b *Buffalo) Main(ctx context.Context, args []string) error {
	cmd.RootCmd.SetArgs(args)
	return cmd.RootCmd.Execute()
}

func (Buffalo) Aliases() []string {
	return []string{"db"}
}

var _ buildcmd.Versioner = &Buffalo{}

func (b *Buffalo) BuildVersion(ctx context.Context, root string) (string, error) {
	return cmd.Version, nil
}

// PkgerDecls() tells Pkger to include Pop related files such as database.yml
// when the buildcmd.Packager interface is called
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

var _ buildcmd.Importer = Buffalo{}

func (Buffalo) BuildImports(ctx context.Context, root string) ([]string, error) {
	dir := filepath.Join(root, "models")
	info, err := here.Dir(dir)
	if err != nil {
		return nil, nil
	}
	return []string{info.ImportPath}, nil
}
