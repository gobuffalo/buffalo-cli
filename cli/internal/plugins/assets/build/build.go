package build

import (
	"context"
	"os"
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/assets/scripts"
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

	c, err := bc.Cmd(ctx, info.Dir, args)
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

func (bc *Builder) tool(ctx context.Context, root string) (string, error) {
	for _, p := range bc.ScopedPlugins() {
		if tp, ok := p.(Tooler); ok {
			return tp.AssetTool(ctx, root)
		}
	}
	return scripts.Tool(root)
}

func (bc *Builder) Cmd(ctx context.Context, root string, args []string) (*exec.Cmd, error) {
	tool := bc.Tool

	var err error
	if len(tool) == 0 {
		tool, err = bc.tool(ctx, root)
		if err != nil {
			return nil, err
		}
	}

	// Fallback on legacy runner
	cmd := plugins.Cmd(ctx, scripts.WebpackBin(root))

	if _, err := scripts.ScriptFor(root, "build"); err == nil {
		cmd = plugins.Cmd(ctx, tool, "run", "build")
	}

	return cmd, nil
}
