package pkger

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here"
	"github.com/markbates/pkger/cmd/pkger/cmds"
	"github.com/markbates/pkger/parser"
)

var _ plugins.Plugin = &Buffalo{}
var _ plugprint.Plugins = &Buffalo{}

const outPath = "pkged.go"

type Buffalo struct {
	OutPath   string
	PluginsFn func() []plugins.Plugin
}

func (b *Buffalo) AfterBuild(ctx context.Context, args []string, err error) error {
	p := b.OutPath
	if len(p) == 0 {
		p = outPath
	}
	os.RemoveAll(p)
	return nil
}

type Decler interface {
	PkgerDecls() (parser.Decls, error)
}

func (b *Buffalo) WithPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if b.PluginsFn != nil {
		plugs = b.PluginsFn()
	}

	var builders []plugins.Plugin
	for _, p := range plugs {
		switch p.(type) {
		case Decler:
			builders = append(builders, p)
		}
	}
	return builders
}

func (b *Buffalo) Build(ctx context.Context, args []string) error {
	return b.Package(ctx, ".")
}

func (b *Buffalo) Package(ctx context.Context, root string) error {
	info, err := here.Current()
	if err != nil {
		return err
	}

	decls, err := parser.Parse(info)
	if err != nil {
		return err
	}

	for _, p := range b.WithPlugins() {
		pd, ok := p.(Decler)
		if !ok {
			continue
		}
		ds, err := pd.PkgerDecls()
		if err != nil {
			return err
		}
		decls = append(decls, ds...)
	}

	os.RemoveAll("pkged.go")
	if err := cmds.Package(info, "pkged.go", decls); err != nil {
		return err
	}

	return nil
}

func (b Buffalo) Name() string {
	return "pkger"
}
