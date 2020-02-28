package newapp

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
)

var _ plugins.Plugin = &Cmd{}
var _ plugcmd.Commander = &Cmd{}
var _ plugcmd.Namer = &Cmd{}

type Cmd struct{}

func (Cmd) PluginName() string {
	return "newapp/cmd"
}

func (Cmd) CmdName() string {
	return "new"
}

func (c *Cmd) Main(ctx context.Context, root string, args []string) error {
	return nil
}
