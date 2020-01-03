package pop

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/plugins/plugprint"
	soda "github.com/gobuffalo/pop/v5/soda/cmd"
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
	soda.RootCmd.SetArgs(args)
	return soda.RootCmd.Execute()
}
