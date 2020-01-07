package pop

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/cli/internal/plugins/pop/internal/soda"
	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
)

type Cmd struct{}

var _ plugins.Plugin = Cmd{}

func (Cmd) Name() string {
	return "pop/cmd"
}

var _ plugprint.NamedCommand = Cmd{}

func (Cmd) CmdName() string {
	return "pop"
}

var _ plugprint.Aliases = Cmd{}

func (Cmd) Aliases() []string {
	return []string{"db"}
}

func (Cmd) Main(ctx context.Context, args []string) error {
	return soda.Main(ctx, args)
}
