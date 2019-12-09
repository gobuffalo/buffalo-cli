package cmdx

import (
	"context"
	"os/exec"
)

func Tidy(ctx context.Context) error {
	c := exec.CommandContext(ctx, "go", "mod", "tidy")
	return c.Run()
}
