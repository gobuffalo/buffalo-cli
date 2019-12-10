package assets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/buffalo-cli/cli/plugins"
	"github.com/gobuffalo/here/there"
	"github.com/gobuffalo/meta/v2"
	"github.com/spf13/pflag"
)

func (b Builder) webpackBin() string {
	s := filepath.Join("node_modules", ".bin", "webpack")
	if runtime.GOOS == "windows" {
		s += ".cmd"
	}
	return s
}

type Builder struct {
	plugins.IO
	Environment string
	// CleanAssets will remove the public/assets folder build compiling
	CleanAssets bool
	// places ./public/assets into ./bin/assets.zip.
	ExtractAssets bool
	SkipAssets    bool
	flagSet       *pflag.FlagSet
}

func (a *Builder) BeforeBuild(ctx context.Context, args []string) error {
	flags := a.PflagSet()
	flags.StringVarP(&a.Environment, "environment", "", "development", "set the environment for the binary")
	flags.Parse(args)

	if a.SkipAssets {
		fmt.Fprintln(a.Stdout(), "skipping assets")
		return nil
	}

	info, err := there.Current()
	if err != nil {
		return err
	}

	app, err := meta.New(info)
	if err != nil {
		return err
	}

	os.Setenv("NODE_ENV", a.Environment)

	tool := "yarnpkg"
	if !app.With["yarn"] {
		tool = "npm"
	}

	type packageJSON struct {
		Scripts map[string]string
	}

	// Fallback on legacy runner
	c := exec.CommandContext(ctx, a.webpackBin())
	scripts := packageJSON{}
	if pf, err := os.Open(filepath.Join(info.Root, "package.json")); err == nil {
		if err = json.NewDecoder(pf).Decode(&scripts); err != nil {
			return err
		}
		if _, ok := scripts.Scripts["build"]; ok {
			c = exec.CommandContext(ctx, tool, "run", "build")
		}
		if err := pf.Close(); err != nil {
			return err
		}
	}

	bb := &bytes.Buffer{}
	c.Stdout = bb
	c.Stderr = bb

	if err := c.Run(); err != nil {
		return err
	}

	if err := a.archive(app); err != nil {
		return err
	}

	return nil
}

func (a Builder) Name() string {
	return "assets builder"
}

func (a Builder) String() string {
	return a.Name()
}

func (a *Builder) BuildFlags() []*pflag.Flag {
	var values []*pflag.Flag
	flags := a.PflagSet()
	flags.VisitAll(func(f *pflag.Flag) {
		values = append(values, f)
	})
	return values
}

func (a *Builder) PflagSet() *pflag.FlagSet {
	if a.flagSet != nil {
		return a.flagSet
	}

	flags := pflag.NewFlagSet(a.String(), pflag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.BoolVar(&a.CleanAssets, "clean-assets", false, "will delete public/assets before calling webpack")
	flags.BoolVarP(&a.ExtractAssets, "extract-assets", "e", false, "extract the assets and put them in a distinct archive")
	flags.BoolVarP(&a.SkipAssets, "skip-assets", "k", false, "skip running webpack and building assets")

	a.flagSet = flags
	return a.flagSet
}

func (a *Builder) PrintFlags(w io.Writer) error {
	flags := a.PflagSet()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
