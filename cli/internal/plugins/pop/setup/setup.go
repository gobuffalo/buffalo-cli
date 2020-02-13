package setup

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/setup"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/pop/v5"
	"github.com/spf13/pflag"
)

var _ plugcmd.Namer = &Setup{}
var _ plugins.Plugin = &Setup{}
var _ setup.BeforeSetuper = &Setup{}
var _ setup.Pflagger = &Setup{}
var _ plugins.Needer = &Setup{}
var _ plugins.Scoper = &Setup{}

type Setup struct {
	flags     *pflag.FlagSet
	dropDB    bool
	pluginsFn plugins.Feeder
}

func (Setup) PluginName() string {
	return "pop/setup"
}

func (Setup) CmdName() string {
	return "pop"
}

func (s *Setup) WithPlugins(f plugins.Feeder) {
	s.pluginsFn = f
}

func (s *Setup) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if s.pluginsFn == nil {
		return plugs
	}

	for _, p := range s.pluginsFn() {
		switch p.(type) {
		case Migrater:
			plugs = append(plugs, p)
		}
	}

	return plugs
}

func (s *Setup) BeforeSetup(ctx context.Context, root string, args []string) error {
	if err := pop.LoadConfigFile(); err != nil {
		return err
	}

	flags := s.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	for _, conn := range pop.Connections {
		if s.dropDB {
			pop.DropDB(conn)
		}
		if err := pop.CreateDB(conn); err != nil {
			return err
		}
	}

	for _, p := range s.ScopedPlugins() {
		switch t := p.(type) {
		case Migrater:
			if err := t.MigrateDB(ctx, root, args); err != nil {
				return err
			}
		}
	}
	return nil
}
