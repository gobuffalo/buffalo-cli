package build

import (
	"context"
	"os"
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/assets/internal/scripts"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"github.com/markbates/safe"
)

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

// BeforeBuild implements the build.BeforeBuilder interface to
// hook into the `buffalo build` lifecycle.
func (a *Builder) BeforeBuild(ctx context.Context, args []string) error {
	return a.Build(ctx, args)
}

// Build implements the build.Builder interface to so it can be run
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

	c, err := bc.cmd(ctx, info.Dir, args)
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

func (bc *Builder) cmd(ctx context.Context, root string, args []string) (*exec.Cmd, error) {
	tool := bc.Tool

	var err error
	if len(tool) == 0 {
		tool, err = scripts.Tool(bc, ctx, root)
		if err != nil {
			return nil, err
		}
	}

	// Fallback on legacy runner
	cmd := plugins.Cmd(ctx, scripts.WebpackBin(root))

	if _, err := scripts.ScriptFor(bc, ctx, root, "build"); err == nil {
		cmd = plugins.Cmd(ctx, tool, "run", "build")
	}

	return cmd, nil
}
