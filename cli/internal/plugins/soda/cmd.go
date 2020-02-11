package soda

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/pop/v5/soda/cmd"
)

var _ plugcmd.Aliaser = Cmd{}
var _ plugcmd.Commander = Cmd{}
var _ plugcmd.Namer = Cmd{}
var _ plugins.Plugin = Cmd{}

type Cmd struct{}

func (Cmd) PluginName() string {
	return "soda/cmd"
}

func (Cmd) CmdName() string {
	return "soda"
}

func (Cmd) CmdAliases() []string {
	return []string{"db", "pop"}
}

func (Cmd) Main(ctx context.Context, root string, args []string) error {
	cmd.RootCmd.SetArgs(args)
	return cmd.RootCmd.Execute()
}
