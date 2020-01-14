package soda

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/plugins"
	"github.com/gobuffalo/buffalo-cli/plugins/plugprint"
)

type Cmd struct{}

var _ plugins.Plugin = Cmd{}

func (Cmd) Name() string {
	return "soda/cmd"
}

var _ plugprint.NamedCommand = Cmd{}

func (Cmd) CmdName() string {
	return "soda"
}

var _ plugprint.Aliases = Cmd{}

func (Cmd) Aliases() []string {
	return []string{"db", "pop"}
}

func (Cmd) Main(ctx context.Context, args []string) error {
	return Main(ctx, args)
}
