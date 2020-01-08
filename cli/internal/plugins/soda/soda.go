package soda

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/pop/v5/soda/cmd"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		Cmd{},
	}
}

func Main(ctx context.Context, args []string) error {
	fmt.Println(">>>TODO cli/internal/plugins/soda/soda.go:18: args ", args)
	cmd.RootCmd.SetArgs(args)
	return cmd.RootCmd.Execute()
}
