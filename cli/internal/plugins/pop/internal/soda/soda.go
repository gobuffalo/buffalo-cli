package soda

import (
	"context"

	"github.com/gobuffalo/pop/v5/soda/cmd"
)

func Main(ctx context.Context, args []string) error {
	cmd.RootCmd.SetArgs(args)
	return cmd.RootCmd.Execute()
}
