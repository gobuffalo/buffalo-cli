package assets

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/assets/webpack"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/here/there"
	"github.com/gobuffalo/meta/v2"
)

type Assets struct {
	stdin         io.Reader
	stdout        io.Writer
	stderr        io.Writer
	Environment   string
	CleanAssets   bool
	ExtractAssets bool
	SkipAssets    bool
	dryRun        bool
}

func (a *Assets) SetStderr(w io.Writer) {
	a.stderr = w
}

func (a *Assets) SetStdin(r io.Reader) {
	a.stdin = r
}

func (a *Assets) SetStdout(w io.Writer) {
	a.stdout = w
}

func (a *Assets) BeforeBuild(ctx context.Context, args []string) error {
	flags := a.flags()
	flags.Parse(args)

	if !a.CleanAssets {
		return nil
	}
	run := genny.WetRunner(ctx)
	if a.dryRun {
		run = genny.DryRunner(ctx)
	}

	info, err := there.Current()
	if err != nil {
		return err
	}

	run.WithRun(func(r *genny.Runner) error {
		r.Delete(filepath.Join(info.Root, "public", "assets"))
		return nil
	})

	run.WithRun(func(r *genny.Runner) error {
		r.Logger.Debugf("setting NODE_ENV = %s", a.Environment)
		return os.Setenv("NODE_ENV", a.Environment)
	})

	app, err := meta.New(info)
	if err != nil {
		return err
	}
	run.WithRun(func(r *genny.Runner) error {
		tool := "yarnpkg"
		if !app.With["yarn"] {
			tool = "npm"
		}

		c := exec.CommandContext(r.Context, tool, "run", "build")
		// if _, err := opts.App.NodeScript("build"); err != nil {
		// Fallback on legacy runner
		c = exec.CommandContext(r.Context, webpack.BinPath)
		// }

		bb := &bytes.Buffer{}
		c.Stdout = bb
		c.Stderr = bb

		if err := r.Exec(c); err != nil {
			r.Logger.Error(bb.String())
			return err
		}
		return nil

	})
	return run.Run()
}

func (a *Assets) AfterBuild(ctx context.Context, args []string) error {
	return nil
}

func (a Assets) Name() string {
	return "assets"
}

func (a Assets) String() string {
	return "assets"
}

func (a *Assets) flags() *cmdx.FlagSet {
	flags := cmdx.NewFlagSet(a.String())
	flags.StringVarP(&a.Environment, "environment", "", "development", "set the environment for the binary")
	flags.BoolVarP(&a.dryRun, "dry-run", "d", false, "dry run")
	flags.BoolVar(&a.CleanAssets, "clean-assets", false, "will delete public/assets before calling webpack")
	flags.BoolVarP(&a.SkipAssets, "skip-assets", "k", false, "skip running webpack and building assets")
	flags.BoolVarP(&a.ExtractAssets, "extract-assets", "e", false, "extract the assets and put them in a distinct archive")
	return flags
}

func (a *Assets) PrintFlags(w io.Writer) error {
	flags := a.flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
