package buildcmd

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gobuffalo/buffalo-cli/plugins"
)

func (bc *BuildCmd) GoBuildArgs() []string {
	buildArgs := []string{"build"}

	if len(bc.Mod) != 0 {
		buildArgs = append(buildArgs, "-mod", bc.Mod)
	}

	buildArgs = append(buildArgs, bc.BuildFlags...)

	if len(bc.Tags) > 0 {
		buildArgs = append(buildArgs, "-tags", bc.Tags)
	}

	bin := bc.Bin
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") {
			bin += ".exe"
		}
		bin = strings.Replace(bin, "/", "\\", -1)
	} else {
		bin = strings.TrimSuffix(bin, ".exe")
	}
	buildArgs = append(buildArgs, "-o", bin)

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
	return buildArgs
}

func (bc *BuildCmd) build(ctx context.Context) error {
	ioe := plugins.CtxIO(ctx)
	cmd := exec.CommandContext(ctx, "go", bc.GoBuildArgs()...)
	cmd.Stdout = ioe.Stdout()
	cmd.Stderr = ioe.Stderr()
	cmd.Stdin = ioe.Stdin()
	fmt.Fprintln(ioe.Stdout(), cmd.Args)
	return cmd.Run()
}
