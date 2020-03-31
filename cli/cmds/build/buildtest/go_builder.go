package buildtest

import (
	"context"
	"os/exec"
)

type GoBuilder func(ctx context.Context, root string, cmd *exec.Cmd) error

func (GoBuilder) PluginName() string {
	return "buildtest/go-builder"
}

func (g GoBuilder) GoBuild(ctx context.Context, root string, cmd *exec.Cmd) error {
	return g(ctx, root, cmd)
}
