package assets

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/buildcmd"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
	"github.com/markbates/safe"
)

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

var _ buildcmd.BeforeBuilder = &Builder{}

// BeforeBuild implements the buildcmd.BeforeBuilder interface to
// hook into the `buffalo build` lifecycle.
func (a *Builder) BeforeBuild(ctx context.Context, args []string) error {
	return a.Build(ctx, args)
}

var _ buildcmd.Builder = &Builder{}

// Build implements the buildcmd.Builder interface to so it can be run
// as `buffalo build assets`.
func (bc *Builder) Build(ctx context.Context, args []string) error {
	var help bool
	flags := bc.Flags()
	flags.StringVarP(&bc.Environment, "environment", "", "development", "set the environment for the binary")
	flags.BoolVarP(&help, "help", "h", false, "print this help")
	flags.Parse(args)

	if help {
		ioe := plugins.CtxIO(ctx)
		return plugprint.Print(ioe.Stdout(), bc)
	}

	if bc.Skip {
		return nil
	}

	os.Setenv("NODE_ENV", bc.Environment)

	info, err := bc.HereInfo()
	if err != nil {
		return err
	}

	c, err := bc.Cmd(info.Dir, ctx, args)
	if err != nil {
		return err
	}

	var fn func() error = c.Run
	for _, p := range bc.ScopedPlugins() {
		if br, ok := p.(AssetBuilder); ok {
			fn = func() error {
				return br.BuildAssets(ctx, c)
			}
			break
		}
	}

	if err := safe.RunE(fn); err != nil {
		return err
	}

	if err := bc.archive(ctx, info.Dir, args); err != nil {
		return err
	}

	return nil
}

func (bc *Builder) Cmd(root string, ctx context.Context, args []string) (*exec.Cmd, error) {
	tool := bc.Tool
	if len(tool) == 0 {
		tool = "npm"
	}

	// Fallback on legacy runner
	c := exec.CommandContext(ctx, bc.webpackBin())

	// parse package.json looking for a custom build script
	scripts := packageJSON{}
	if pf, err := os.Open(filepath.Join(root, "package.json")); err == nil {
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
