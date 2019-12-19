package assets

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/here/there"
	"github.com/gobuffalo/meta/v2"
)

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

// BeforeBuild implements the buildcmd.BeforeBuilder interface to
// hook into the `buffalo build` lifecycle.
func (a *Builder) BeforeBuild(ctx context.Context, args []string) error {
	return a.Build(ctx, args)
}

// Build implements the buildcmd.Builder interface to so it can be run
// as `buffalo build assets`.
func (bc *Builder) Build(ctx context.Context, args []string) error {
	flags := bc.Flags()
	flags.StringVarP(&bc.Environment, "environment", "", "development", "set the environment for the binary")
	flags.Parse(args)

	if bc.SkipAssets {
		return nil
	}

	os.Setenv("NODE_ENV", bc.Environment)

	info, err := there.Current()
	if err != nil {
		return err
	}

	if i, ok := ctx.Value("here.Current").(here.Info); ok {
		info = i
	}

	app, err := meta.New(info)
	if err != nil {
		return err
	}

	c, err := bc.Cmd(app, ctx, args)
	if err != nil {
		return err
	}

	if err := c.Run(); err != nil {
		return err
	}

	if err := bc.archive(ctx, args); err != nil {
		return err
	}

	return nil
}

func (bc *Builder) Cmd(app *meta.App, ctx context.Context, args []string) (*exec.Cmd, error) {
	tool := "yarnpkg"
	if !app.With["yarn"] {
		tool = "npm"
	}

	// Fallback on legacy runner
	c := exec.CommandContext(ctx, bc.webpackBin())

	// parse package.json looking for a custom build script
	scripts := packageJSON{}
	if pf, err := os.Open(filepath.Join(app.Info.Root, "package.json")); err == nil {
		if err = json.NewDecoder(pf).Decode(&scripts); err != nil {
			return nil, err
		}
		if _, ok := scripts.Scripts["build"]; ok {
			c = exec.CommandContext(ctx, tool, "run", "build")
		}
		if err := pf.Close(); err != nil {
			return nil, err
		}
	}

	ioe := plugins.CtxIO(ctx)
	c.Stdout = ioe.Stdout()
	c.Stderr = ioe.Stderr()
	c.Stdin = ioe.Stdin()
	return c, nil
}
