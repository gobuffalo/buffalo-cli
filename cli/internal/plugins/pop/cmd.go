package pop

import (
	"github.com/gobuffalo/buffalo-cli/internal/plugins"
	"github.com/gobuffalo/buffalo-cli/internal/plugins/plugprint"
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
