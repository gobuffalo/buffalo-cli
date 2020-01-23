package soda

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
	"github.com/gobuffalo/pop/v5/soda/cmd"
)

var _ plugins.Plugin = Cmd{}
var _ plugprint.Aliases = Cmd{}
var _ plugprint.NamedCommand = Cmd{}

type Cmd struct{}

func (Cmd) Name() string {
	return "soda/cmd"
}

func (Cmd) CmdName() string {
	return "soda"
}

func (Cmd) Aliases() []string {
	return []string{"db", "pop"}
}

func (Cmd) Main(ctx context.Context, root string, args []string) error {
	cmd.RootCmd.SetArgs(args)
	return cmd.RootCmd.Execute()
}
