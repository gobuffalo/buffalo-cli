package buildcmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/build"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/gobuffalo/here/there"
	"github.com/gobuffalo/meta"
)

func (bc *BuildCmd) Main(ctx context.Context, args []string) error {
	info, err := there.Current()
	if err != nil {
		return err
	}

	opts := &build.Options{
		App: meta.New(info.Root),
	}

	flags := bc.flagSet(opts)
	if err = flags.Parse(args); err != nil {
		return err
	}

	if bc.help {
		return plugprint.Print(bc.Stdout(), bc)
	}

	plugs := bc.WithPlugins()

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

	builders := bc.WithPlugins()
	for _, p := range builders {
		if bb, ok := p.(BeforeBuilder); ok {
			plugins.SetIO(bc, p)
			if err := bb.BeforeBuild(ctx, args); err != nil {
				return err
			}
		}
	}

	if bc.verbose {
		bc.BuildFlags = append(bc.BuildFlags, "-v")
	}

	if !bc.skipTemplateValidation {
		for _, p := range plugs {
			tv, ok := p.(TemplatesValidator)
			if !ok {
				continue
			}
			plugins.SetIO(bc, p)
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
		plugins.SetIO(bc, p)
		if err := pkg.Package(ctx, info.Root); err != nil {
			return err
		}
	}

	version := time.Now().Format(time.RFC3339)
	for _, p := range plugs {
		bv, ok := p.(BuildVersioner)
		if !ok {
			continue
		}
		plugins.SetIO(bc, p)
		s, err := bv.BuildVersion(ctx, info.Root)
		if err != nil {
			return err
		}
		if len(s) == 0 {
			continue
		}
		version = strings.TrimSpace(s)
	}
	fmt.Println("version: ", version)

	if err := bc.build(ctx); err != nil {
		return err
	}

	for _, p := range builders {
		if bb, ok := p.(AfterBuilder); ok {
			plugins.SetIO(bc, p)
			if err := bb.AfterBuild(ctx, args); err != nil {
				return err
			}
		}
	}
	return nil
}
