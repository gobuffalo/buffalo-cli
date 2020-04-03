package build

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/markbates/safe"
)

func (bc *Cmd) GoCmd(ctx context.Context, root string) (*exec.Cmd, error) {
	buildArgs := []string{"build"}

	info, err := here.Dir(root)
	if err != nil {
		return nil, err
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
	buildArgs = append(buildArgs, "-o", bin)

	if len(bc.mod) != 0 {
		buildArgs = append(buildArgs, "-mod", bc.mod)
	}

	buildArgs = append(buildArgs, bc.buildFlags...)

	tags, err := bc.buildTags(ctx, root)
	if err != nil {
		return nil, err
	}

	if len(tags) > 0 {
		buildArgs = append(buildArgs, "-tags", strings.Join(tags, " "))
	}

	flags := []string{}

	if bc.static {
		flags = append(flags, "-linkmode external", "-extldflags \"-static\"")
	}

	// Add any additional ldflags passed in to the build args
	if len(bc.ldFlags) > 0 {
		flags = append(flags, bc.ldFlags)
	}
	if len(flags) > 0 {
		buildArgs = append(buildArgs, "-ldflags", strings.Join(flags, " "))
	}

	cmd := exec.CommandContext(ctx, "go", buildArgs...)

	plugs := bc.ScopedPlugins()
	cmd.Stdout = plugio.Stdout(plugs...)
	cmd.Stderr = plugio.Stderr(plugs...)
	cmd.Stdin = plugio.Stdin(plugs...)

	return cmd, nil
}

func (cmd *Cmd) buildTags(ctx context.Context, root string) ([]string, error) {
	var tags []string
	if len(cmd.tags) > 0 {
		tags = append(tags, cmd.tags)
	}

	for _, p := range cmd.ScopedPlugins() {
		t, ok := p.(Tagger)
		if !ok {
			continue
		}
		bt, err := t.BuildTags(ctx, root)
		if err != nil {
			return nil, err
		}
		tags = append(tags, bt...)
	}

	return tags, nil
}

func (bc *Cmd) build(ctx context.Context, root string, args []string) error {
	cmd, err := bc.GoCmd(ctx, root)
	if err != nil {
		return err
	}
	fmt.Println(cmd.Args)

	for _, p := range bc.ScopedPlugins() {
		if br, ok := p.(Runner); ok {
			return safe.RunE(func() error {
				return br.RunBuild(ctx, cmd)
			})
		}
	}

	return cmd.Run()
}
