package pkger

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/here"
	"github.com/markbates/pkger/cmd/pkger/cmds"
	"github.com/markbates/pkger/parser"
)

const outPath = "pkged.go"

var _ build.AfterBuilder = &Builder{}
var _ build.Builder = &Builder{}
var _ build.Packager = &Builder{}
var _ plugins.Plugin = &Builder{}
var _ plugins.PluginNeeder = &Builder{}
var _ plugins.PluginScoper = &Builder{}

type Builder struct {
	OutPath   string
	pluginsFn plugins.PluginFeeder
}

func (b *Builder) WithPlugins(f plugins.PluginFeeder) {
	b.pluginsFn = f
}

func (b *Builder) AfterBuild(ctx context.Context, root string, args []string, err error) error {
	p := b.OutPath
	if len(p) == 0 {
		p = outPath
	}
	os.RemoveAll(p)
	return nil
}

func (b *Builder) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if b.pluginsFn != nil {
		plugs = b.pluginsFn()
	}

	return plugs
}

func (b *Builder) Build(ctx context.Context, root string, args []string) error {
	return b.Package(ctx, root, nil)
}

func (b *Builder) Package(ctx context.Context, root string, files []string) error {
	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	decls, err := parser.Parse(info)
	if err != nil {
		return err
	}
	for _, f := range files {
		d, err := parser.NewInclude(info, f)
		if err != nil {
			return err
		}
		decls = append(decls, d)
	}

	os.RemoveAll("pkged.go")
	if err := cmds.Package(info, "pkged.go", decls); err != nil {
		return err
	}

	return nil
}

func (b Builder) Name() string {
	return "pkger"
}
