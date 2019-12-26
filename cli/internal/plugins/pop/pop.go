package pop

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/pop/v5/soda/cmd"
	"github.com/markbates/pkger"
)

const filePath = "/database.yml"

type Buffalo struct{}

var _ plugins.Plugin = Buffalo{}

func (Buffalo) Name() string {
	return "buffalo-pop"
}

var _ plugprint.NamedCommand = Buffalo{}

func (Buffalo) CmdName() string {
	return "pop"
}

// Main adds the `buffalo pop` sub-command.
func (b *Buffalo) Main(ctx context.Context, args []string) error {
	cmd.RootCmd.SetArgs(args)
	return cmd.RootCmd.Execute()
}

var _ plugprint.Aliases = Buffalo{}

func (Buffalo) Aliases() []string {
	return []string{"db"}
}

var _ buildcmd.Versioner = &Buffalo{}

func (b *Buffalo) BuildVersion(ctx context.Context, root string) (string, error) {
	return cmd.Version, nil
}

var _ buildcmd.PackFiler = &Buffalo{}

func (b *Buffalo) PackageFiles(ctx context.Context, root string) ([]string, error) {
	return []string{
		filepath.Join(root, filePath),
	}, nil
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
