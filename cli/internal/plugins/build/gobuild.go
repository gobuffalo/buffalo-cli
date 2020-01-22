package build

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/markbates/safe"
)

func (bc *Cmd) GoCmd(ctx context.Context) (*exec.Cmd, error) {
	buildArgs := []string{"build"}

	info, err := bc.HereInfo()
	if err != nil {
		return nil, err
	}

	bin := bc.Bin
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
	buildArgs = append(buildArgs, "-o", bin)

	if len(bc.Mod) != 0 {
		buildArgs = append(buildArgs, "-mod", bc.Mod)
	}

	buildArgs = append(buildArgs, bc.BuildFlags...)

	if len(bc.Tags) > 0 {
		buildArgs = append(buildArgs, "-tags", bc.Tags)
	}

	flags := []string{}

	if bc.Static {
		flags = append(flags, "-linkmode external", "-extldflags \"-static\"")
	}

	// Add any additional ldflags passed in to the build args
	if len(bc.LDFlags) > 0 {
		flags = append(flags, bc.LDFlags)
	}
	if len(flags) > 0 {
		buildArgs = append(buildArgs, "-ldflags", strings.Join(flags, " "))
	}

	cmd := exec.CommandContext(ctx, "go", buildArgs...)

	ioe := plugins.CtxIO(ctx)
	cmd.Stdout = ioe.Stdout()
	cmd.Stderr = ioe.Stderr()
	cmd.Stdin = ioe.Stdin()

	if testing.Verbose() {
		fmt.Fprintln(ioe.Stdout(), cmd.Args)
	}

	return cmd, nil
}

func (bc *Cmd) build(ctx context.Context, args []string) error {
	cmd, err := bc.GoCmd(ctx)
	if err != nil {
		return err
	}

	for _, p := range bc.ScopedPlugins() {
		if br, ok := p.(Runner); ok {
			return safe.RunE(func() error {
				return br.RunBuild(ctx, cmd)
			})
		}
	}

	return cmd.Run()
}
