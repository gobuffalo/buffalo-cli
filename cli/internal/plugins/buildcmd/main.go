package buildcmd

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/plugins/plugprint"
)

func (bc *BuildCmd) beforeBuild(ctx context.Context, args []string) error {
	builders := bc.ScopedPlugins()
	for _, p := range builders {
		if bb, ok := p.(BeforeBuilder); ok {
			if err := bb.BeforeBuild(ctx, args); err != nil {
				return err
			}
		}
	}
	return nil
}

func (bc *BuildCmd) afterBuild(ctx context.Context, args []string, err error) error {
	builders := bc.ScopedPlugins()
	for _, p := range builders {
		if bb, ok := p.(AfterBuilder); ok {
			if err := bb.AfterBuild(ctx, args, err); err != nil {
				return err
			}
		}
	}
	return err
}

func (bc *BuildCmd) Main(ctx context.Context, args []string) error {
	flags := bc.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}
	if bc.verbose {
		bc.BuildFlags = append(bc.BuildFlags, "-v")
	}

	ioe := plugins.CtxIO(ctx)

	info, err := bc.HereInfo()
	if err != nil {
		return err
	}

	plugs := bc.ScopedPlugins()

	if len(flags.Args()) > 0 {
		n := flags.Args()[0]
		for _, p := range plugs {
			b, ok := p.(Builder)
			if !ok {
				continue
			}
			if p.Name() == n {
				return b.Build(ctx, args[1:])
			}
		}
		return fmt.Errorf("unknown command %q", n)
	}

	if bc.help {
		return plugprint.Print(ioe.Stdout(), bc)
	}

	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%s", e)
			}
			bc.afterBuild(ctx, args, err)
		}
	}()

	if err = bc.beforeBuild(ctx, args); err != nil {
		return bc.afterBuild(ctx, args, err)
	}

	if !bc.skipTemplateValidation {
		for _, p := range plugs {
			tv, ok := p.(TemplatesValidator)
			if !ok {
				continue
			}
			if err = tv.ValidateTemplates(info.Root); err != nil {
				return bc.afterBuild(ctx, args, err)
			}
		}
	}

	if err := bc.pack(ctx, info, plugs); err != nil {
		return bc.afterBuild(ctx, args, err)
	}

	err = bc.build(ctx) // go build ...
	return bc.afterBuild(ctx, args, err)
}
