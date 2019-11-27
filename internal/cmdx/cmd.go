package cmdx

import (
	"context"
	"os/exec"
)

func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	c := exec.CommandContext(ctx, name, arg...)
	c.Stdin = Stdin(ctx)
	c.Stdout = Stdout(ctx)
	c.Stderr = Stderr(ctx)
	return c
}
