package cli

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/internal/v1/cmd"
)

func Main(ctx context.Context, args []string) error {
	c := cmd.RootCmd
	c.SetArgs(args)
	return c.Execute()
}
