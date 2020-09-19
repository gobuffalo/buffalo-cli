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

// Cmd is a Fixer in charge of generating cmd/buffalo/main.go
// so the buffalo CLI can be configured.
type Cmd struct{}

func (*Cmd) PluginName() string {
	return "cli/fixer"
}

func (*Cmd) CmdName() string {
	return "cli"
}

// Fix attempts to generate cmd/buffalo/main.go if it does not exist
// otherwise nothing happens.
func (fixer *Cmd) Fix(ctx context.Context, root string, args []string) error {
	g := &cligen.Generator{}
	err := g.Generate(ctx, root, args)
	if err == cligen.ErrCLIMainExists {
		return nil
	}

	return err
}
