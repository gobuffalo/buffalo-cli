package builder

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/pop/v5/soda/cmd"
)

var _ build.Importer = Builder{}
var _ build.PackFiler = &Builder{}
var _ build.Tagger = &Builder{}
var _ build.Versioner = &Builder{}
var _ plugins.Plugin = Builder{}

const filePath = "/database.yml"

type Builder struct{}

func (Builder) PluginName() string {
	return "pop/builder"
}

func (bd *Builder) BuildTags(ctx context.Context, root string) ([]string, error) {
	var args []string
	dy := filepath.Join(root, "database.yml")
	if _, err := os.Stat(dy); err != nil {
		return args, nil
	}

	b, err := ioutil.ReadFile(dy)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(b, []byte("sqlite")) {
		args = append(args, "sqlite")
	}
	return args, nil
}

func (b *Builder) BuildVersion(ctx context.Context, root string) (string, error) {
	return cmd.Version, nil
}

func (b *Builder) PackageFiles(ctx context.Context, root string) ([]string, error) {
	return []string{
		filepath.Join(root, filePath),
	}, nil
}

func (Builder) BuildImports(ctx context.Context, root string) ([]string, error) {
	dir := filepath.Join(root, "models")
	info, err := here.Dir(dir)
	if err != nil {
		return nil, nil
	}
	return []string{info.ImportPath}, nil
}
