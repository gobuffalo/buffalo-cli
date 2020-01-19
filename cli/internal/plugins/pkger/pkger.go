package pkger

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here"
	"github.com/markbates/pkger/cmd/pkger/cmds"
	"github.com/markbates/pkger/parser"
)

const outPath = "pkged.go"

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&Builder{},
	}
}

type Builder struct {
	OutPath   string
	pluginsFn plugins.PluginFeeder
}

var _ plugins.PluginNeeder = &Builder{}

func (b *Builder) WithPlugins(f plugins.PluginFeeder) {
	b.pluginsFn = f
}

var _ build.AfterBuilder = &Builder{}

func (b *Builder) AfterBuild(ctx context.Context, args []string, err error) error {
	p := b.OutPath
	if len(p) == 0 {
		p = outPath
	}
	os.RemoveAll(p)
	return nil
}

var _ plugins.PluginScoper = &Builder{}

func (b *Builder) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if b.pluginsFn != nil {
		plugs = b.pluginsFn()
	}

	return plugs
}

var _ build.Builder = &Builder{}

func (b *Builder) Build(ctx context.Context, args []string) error {
	return b.Package(ctx, ".", nil)
}

var _ build.Packager = &Builder{}

func (b *Builder) Package(ctx context.Context, root string, files []string) error {
	info, err := here.Current()
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

var _ plugins.Plugin = &Builder{}

func (b Builder) Name() string {
	return "pkger"
}
