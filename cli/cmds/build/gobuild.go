package build

import (
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
)

func (bc *Cmd) buildArgs(ctx context.Context, root string) ([]string, error) {
	args := []string{"build"}

	info, err := here.Dir(root)
	if err != nil {
		return nil, plugins.Wrap(bc, err)
	}

	bin := bc.bin
	if len(bin) == 0 {
		bin = filepath.Join("bin", info.Name)
	}

	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") {
			bin += ".exe"
		}
		bin = strings.Replace(bin, "/", "\\", -1)
	} else {
		bin = strings.TrimSuffix(bin, ".exe")
	}
	args = AppendArg(args, "-o", bin)

	if len(bc.mod) != 0 {
		args = AppendArg(args, "-mod", bc.mod)
	}

	args = append(args, bc.buildFlags...)

	flags := []string{}
	if bc.static {
		flags = append(flags, "-linkmode external", "-extldflags \"-static\"")
	}

	// Add any additional ldflags passed in to the build args
	if len(bc.ldFlags) > 0 {
		flags = append(flags, bc.ldFlags)
	}
	if len(flags) > 0 {
		args = AppendArg(args, "-ldflags", flags...)
	}

	if len(bc.tags) > 0 {
		args = AppendArg(args, "-tags", bc.tags)
	}

	for _, p := range bc.ScopedPlugins() {
		if bt, ok := p.(BuildArger); ok {
			ags, err := bt.GoBuildArgs(ctx, root, args)
			if err != nil {
				return nil, err
			}

			switch len(ags) {
			case 0:
				continue
			case 1:
				args = AppendArg(args, ags[0])
			default:
				args = AppendArg(args, ags[0], ags[1:]...)
			}

		}
	}

	return args, nil
}

func (bc *Cmd) build(ctx context.Context, root string) error {
	buildArgs, err := bc.buildArgs(ctx, root)
	if err != nil {
		return plugins.Wrap(bc, err)
	}

	plugs := bc.ScopedPlugins()
	for _, p := range plugs {
		if br, ok := p.(GoBuilder); ok {
			if err := br.GoBuild(ctx, root, buildArgs); err != nil {
				return plugins.Wrap(br, err)
			}
		}
	}

	cmd := exec.CommandContext(ctx, "go", buildArgs...)
	cmd.Stdin = plugio.Stdin(plugs...)
	cmd.Stdout = plugio.Stdout(plugs...)
	cmd.Stderr = plugio.Stderr(plugs...)
	return cmd.Run()
}
