package build

import (
	"context"
	"fmt"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugfind"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/markbates/safe"
)

func byBuilder(f plugfind.Finder) plugfind.Finder {
	fn := func(name string, plugs []plugins.Plugin) plugins.Plugin {
		p := f.Find(name, plugs)
		if p == nil {
			return nil
		}
		if c, ok := p.(Builder); ok {
			if c.PluginName() == name {
				return p
			}
		}
		return nil
	}
	return plugfind.FinderFn(fn)
}

func (bc *Cmd) Main(ctx context.Context, root string, args []string) error {
	flags := bc.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	if len(flags.Args()) == 0 && bc.help {
		return plugprint.Print(plugio.Stdout(bc.ScopedPlugins()...), bc)
	}

	type builder func(ctx context.Context, root string, args []string) error

	var build builder = bc.build

	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	plugs := bc.ScopedPlugins()

	if len(flags.Args()) > 0 {
		name := flags.Args()[0]
		fn := plugfind.Background()
		fn = byBuilder(fn)
		fn = plugcmd.ByNamer(fn)
		fn = plugcmd.ByAliaser(fn)

		p := fn.Find(name, plugs)
		if p == nil {
			return fmt.Errorf("unknown builder %q", name)
		}
		b, ok := p.(Builder)
		if !ok {
			return fmt.Errorf("unknown builder %q", name)
		}
		build = b.Build
		args = args[1:]
	}

	if err = bc.beforeBuild(ctx, root, args); err != nil {
		return bc.afterBuild(ctx, root, args, err)
	}

	if !bc.SkipTemplateValidation {
		for _, p := range plugs {
			tv, ok := p.(TemplatesValidator)
			if !ok {
				continue
			}
			err = safe.RunE(func() error {
				return tv.ValidateTemplates(info.Dir)
			})
			if err != nil {
				return bc.afterBuild(ctx, root, args, err)
			}
		}
	}

	err = safe.RunE(func() error {
		return bc.pack(ctx, info, plugs)
	})
	if err != nil {
		return bc.afterBuild(ctx, root, args, err)
	}

	err = safe.RunE(func() error {
		return build(ctx, root, args)
	})

	return bc.afterBuild(ctx, root, args, err)

}
func (bc *Cmd) beforeBuild(ctx context.Context, root string, args []string) error {
	builders := bc.ScopedPlugins()
	for _, p := range builders {
		if bb, ok := p.(BeforeBuilder); ok {
			err := safe.RunE(func() error {
				return bb.BeforeBuild(ctx, root, args)
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (bc *Cmd) afterBuild(ctx context.Context, root string, args []string, err error) error {
	builders := bc.ScopedPlugins()
	for _, p := range builders {
		if bb, ok := p.(AfterBuilder); ok {
			err := safe.RunE(func() error {
				return bb.AfterBuild(ctx, root, args, err)
			})
			if err != nil {
				return err
			}
		}
	}
	return err
}
