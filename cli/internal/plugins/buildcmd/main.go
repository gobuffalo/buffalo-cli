package buildcmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here/there"
)

func (bc *BuildCmd) beforeBuild(ctx context.Context, args []string) error {
	builders := bc.Plugins()
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
	builders := bc.Plugins()
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
	flags := bc.flagSet()
	if err := flags.Parse(args); err != nil {
		return err
	}
	if bc.verbose {
		bc.BuildFlags = append(bc.BuildFlags, "-v")
	}

	ioe := plugins.CtxIO(ctx)
	if bc.help {
		return plugprint.Print(ioe.Stdout(), bc)
	}

	info, err := there.Current()
	if err != nil {
		return err
	}

	defer func() {
		var err error
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%s", e)
			}
		}
		bc.afterBuild(ctx, args, err)
	}()

	plugs := bc.Plugins()

	if len(flags.Args()) > 0 {
		n := flags.Args()[0]
		for _, p := range plugs {
			b, ok := p.(Builder)
			if !ok {
				continue
			}
			if p.Name() == n {
				return b.Build(ctx, args)
			}
		}
		return fmt.Errorf("unknown command %q", n)
	}

	for _, p := range plugs {
		if bb, ok := p.(BeforeBuilder); ok {
			if err := bb.BeforeBuild(ctx, args); err != nil {
				return err
			}
		}
	}

	if !bc.skipTemplateValidation {
		for _, p := range plugs {
			tv, ok := p.(TemplatesValidator)
			if !ok {
				continue
			}
			if err := tv.ValidateTemplates(filepath.Join(info.Root, "templates")); err != nil {
				return err
			}
		}
	}

	for _, p := range plugs {
		pkg, ok := p.(Packager)
		if !ok {
			continue
		}
		if err := pkg.Package(ctx, info.Root); err != nil {
			return err
		}
	}

	return bc.build(ctx) // go build ...
}
