package pkger

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here"
	"github.com/markbates/pkger/cmd/pkger/cmds"
	"github.com/markbates/pkger/parser"
)

const outPath = "pkged.go"

type Buffalo struct {
	OutPath   string
	pluginsFn plugins.PluginFeeder
}

var _ plugins.PluginNeeder = &Buffalo{}

func (b *Buffalo) WithPlugins(f plugins.PluginFeeder) {
	b.pluginsFn = f
}

var _ buildcmd.AfterBuilder = &Buffalo{}

func (b *Buffalo) AfterBuild(ctx context.Context, args []string, err error) error {
	p := b.OutPath
	if len(p) == 0 {
		p = outPath
	}
	os.RemoveAll(p)
	return nil
}

var _ plugins.PluginScoper = &Buffalo{}

func (b *Buffalo) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if b.pluginsFn != nil {
		plugs = b.pluginsFn()
	}

	return plugs
}

var _ buildcmd.Builder = &Buffalo{}

func (b *Buffalo) Build(ctx context.Context, args []string) error {
	return b.Package(ctx, ".", nil)
}

var _ buildcmd.Packager = &Buffalo{}

func (b *Buffalo) Package(ctx context.Context, root string, files []string) error {
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

var _ plugins.Plugin = &Buffalo{}

func (b Buffalo) Name() string {
	return "pkger"
}
