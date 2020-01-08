package soda

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/pop/v5/soda/cmd"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		Cmd{},
	}
}

func Main(ctx context.Context, args []string) error {
	cmd.RootCmd.SetArgs(args)
	return cmd.RootCmd.Execute()
}
