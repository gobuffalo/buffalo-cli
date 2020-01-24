package grifts

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	grifts "github.com/markbates/grift/cmd"
)

var _ plugins.Plugin = &Cmd{}
var _ plugins.Aliases = &Cmd{}

type Cmd struct{}

func (Cmd) Name() string {
	return "grifts"
}

func (Cmd) Aliases() []string {
	return []string{"task", "tasks", "t"}
}

func (cmd *Cmd) Main(ctx context.Context, root string, args []string) error {
	return grifts.Run("buffalo grifts", args)
}
