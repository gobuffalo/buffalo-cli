package cmd

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	grifts "github.com/markbates/grift/cmd"
)

var _ plugcmd.Aliaser = &Cmd{}
var _ plugcmd.Namer = &Cmd{}
var _ plugins.Plugin = &Cmd{}

type Cmd struct{}

func (Cmd) PluginName() string {
	return "grifts/cmd"
}

func (Cmd) CmdName() string {
	return "grift"
}

func (Cmd) CmdAliases() []string {
	return []string{"task", "tasks", "t"}
}

func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	return grifts.Run("buffalo grifts", args)
}
