package build

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"github.com/markbates/safe"
)

func (bc *Cmd) beforeBuild(ctx context.Context, args []string) error {
	builders := bc.ScopedPlugins()
	for _, p := range builders {
		if bb, ok := p.(BeforeBuilder); ok {
			err := safe.RunE(func() error {
				return bb.BeforeBuild(ctx, args)
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (bc *Cmd) afterBuild(ctx context.Context, args []string, err error) error {
	builders := bc.ScopedPlugins()
	for _, p := range builders {
		if bb, ok := p.(AfterBuilder); ok {
			err := safe.RunE(func() error {
				return bb.AfterBuild(ctx, args, err)
			})
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (bc *Cmd) Main(ctx context.Context, root string, args []string) error {
	flags := bc.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	ioe := plugins.CtxIO(ctx)

	if len(flags.Args()) == 0 && bc.help {
		return plugprint.Print(ioe.Stdout(), bc)
	}

	type builder func(ctx context.Context, args []string) error

	var build builder = bc.build

	info, err := bc.HereInfo()
	if err != nil {
		return err
	}

	plugs := bc.ScopedPlugins()

	if len(flags.Args()) > 0 {
		n := flags.Args()[0]
		cmds := plugins.Commands(plugs)
		p, err := cmds.Find(n)
		if err != nil {
			return err
		}
		b, ok := p.(Builder)
		if !ok {
			return fmt.Errorf("unknown command %q", n)
		}
		build = b.Build
		args = args[1:]
	}

	if err = bc.beforeBuild(ctx, args); err != nil {
		return bc.afterBuild(ctx, args, err)
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
				return bc.afterBuild(ctx, args, err)
			}
		}
	}

	err = safe.RunE(func() error {
		return bc.pack(ctx, info, plugs)
	})
	if err != nil {
		return bc.afterBuild(ctx, args, err)
	}

	err = safe.RunE(func() error {
		return build(ctx, args)
	})

	return bc.afterBuild(ctx, args, err)
}
