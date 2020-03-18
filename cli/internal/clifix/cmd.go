package clifix

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/fix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/cligen"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
)

var _ plugins.Plugin = &Cmd{}
var _ plugcmd.Namer = &Cmd{}
var _ fix.Fixer = &Cmd{}

type Cmd struct {
}

func (*Cmd) PluginName() string {
	return "cli/fixer"
}

func (*Cmd) CmdName() string {
	return "cli"
}

func (fixer *Cmd) Fix(ctx context.Context, root string, args []string) error {
	g := &cligen.Generator{}
	return g.Generate(ctx, root, args)
}
