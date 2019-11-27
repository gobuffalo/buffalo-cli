package cmdx

import (
	"context"
)

func Tidy(ctx context.Context) error {
	c := CommandContext(ctx, "go", "mod", "tidy")
	return c.Run()
}
