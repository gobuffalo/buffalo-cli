package bzr

import (
	"context"
	"os/exec"
)

type runner func(ctx context.Context, root string, cmd *exec.Cmd) error

func (runner) PluginName() string {
	return "bzr-runner"
}

func (r runner) RunBzr(ctx context.Context, root string, cmd *exec.Cmd) error {
	return r(ctx, root, cmd)
}
