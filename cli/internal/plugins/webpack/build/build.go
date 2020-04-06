package build

import (
	"context"
	"os"
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack/internal/scripts"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
)

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

// BeforeBuild implements the build.BeforeBuilder interface to
// hook into the `buffalo build` lifecycle.
func (a *Builder) BeforeBuild(ctx context.Context, root string, args []string) error {
	return a.Build(ctx, root, args)
}

// Build implements the build.Builder interface to so it can be run
// as `buffalo build webpack`.
func (bc *Builder) Build(ctx context.Context, root string, args []string) error {
	var help bool
	flags := bc.Flags()
	flags.StringVarP(&bc.environment, "environment", "", "development", "set the environment for the binary")
	flags.BoolVarP(&help, "help", "h", false, "print this help")
	flags.Parse(args)

	if help {
		return plugprint.Print(plugio.Stdout(bc.ScopedPlugins()...), bc)
	}

	if bc.skip {
		return nil
	}

	ne := os.Getenv("NODE_ENV")
	defer os.Setenv("NODE_ENV", ne)
	os.Setenv("NODE_ENV", bc.environment)

	c, err := bc.cmd(ctx, root, args)
	if err != nil {
		return plugins.Wrap(bc, err)
	}

	var fn func() error = c.Run
	for _, p := range bc.ScopedPlugins() {
		if br, ok := p.(AssetBuilder); ok {
			fn = func() error {
				return br.BuildAssets(ctx, root, c)
			}
			break
		}
	}

	if err := fn(); err != nil {
		return plugins.Wrap(bc, err)
	}

	if err := bc.archive(ctx, root, args); err != nil {
		return plugins.Wrap(bc, err)
	}

	return nil
}

func (bc *Builder) cmd(ctx context.Context, root string, args []string) (*exec.Cmd, error) {
	tool := bc.tool

	var err error
	if len(tool) == 0 {
		tool, err = scripts.Tool(bc, ctx, root)
		if err != nil {
			return nil, plugins.Wrap(bc, err)
		}
	}

	// Fallback on legacy runner
	cmd := exec.CommandContext(ctx, scripts.WebpackBin(root))

	if _, err := scripts.ScriptFor(bc, ctx, root, "build"); err == nil {
		cmd = exec.CommandContext(ctx, tool, "run", "build")
	}

	plugs := bc.ScopedPlugins()
	cmd.Stdin = plugio.Stdin(plugs...)
	cmd.Stdout = plugio.Stdout(plugs...)
	cmd.Stderr = plugio.Stderr(plugs...)

	return cmd, nil
}
