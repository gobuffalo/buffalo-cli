package test

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/setup"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
)

var _ plugcmd.Namer = &Setup{}
var _ plugins.Needer = &Setup{}
var _ plugins.Plugin = &Setup{}
var _ setup.Setuper = &Setup{}

type Setup struct {
	pluginsFn plugins.Feeder
}

func (s Setup) PluginName() string {
	return "test/setup"
}

func (s Setup) CmdName() string {
	return "test"
}

func (s *Setup) WithPlugins(f plugins.Feeder) {
	s.pluginsFn = f
}

func (s Setup) Setup(ctx context.Context, root string, args []string) error {
	tc := &Cmd{}
	if s.pluginsFn != nil {
		for _, p := range s.pluginsFn() {
			if t, ok := p.(*Cmd); ok {
				tc = t
				break
			}
		}
	}
	tc.WithPlugins(s.pluginsFn)
	return tc.Main(ctx, root, args)
}
