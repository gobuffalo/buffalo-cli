package setup

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/setup"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/gobuffalo/pop/v5"
	"github.com/spf13/pflag"
)

var _ plugcmd.Namer = &Setup{}
var _ plugins.Needer = &Setup{}
var _ plugins.Plugin = &Setup{}
var _ plugins.Scoper = &Setup{}
var _ setup.Pflagger = &Setup{}
var _ setup.Setuper = &Setup{}

type Setup struct {
	flags     *pflag.FlagSet
	dropDB    bool
	help      bool
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
		case DBSeeder:
			plugs = append(plugs, p)
		}
	}

	return plugs
}

func (s *Setup) Setup(ctx context.Context, root string, args []string) error {
	if err := pop.LoadConfigFile(); err != nil {
		return plugins.Wrap(s, err)
	}

	flags := s.Flags()
	if err := flags.Parse(args); err != nil {
		return plugins.Wrap(s, err)
	}

	plugs := s.ScopedPlugins()

	if s.help {
		return plugprint.Print(plugio.Stdout(plugs...), s)
	}

	for _, conn := range pop.Connections {
		if s.dropDB {
			pop.DropDB(conn)
		}
		if err := pop.CreateDB(conn); err != nil {
			return plugins.Wrap(s, err)
		}
	}

	for _, p := range plugs {
		if t, ok := p.(Migrater); ok {
			if err := t.MigrateDB(ctx, root, args); err != nil {
<<<<<<< HEAD
				return plugins.Wrap(p, err)
=======
				return plugins.Wrap(s, err)
>>>>>>> i was only asking
			}
		}
	}

	for _, p := range plugs {
		if t, ok := p.(DBSeeder); ok {
			if err := t.SeedDB(ctx, root, args); err != nil {
<<<<<<< HEAD
				return plugins.Wrap(p, err)
=======
				return plugins.Wrap(s, err)
>>>>>>> i was only asking
			}
		}
	}

	return nil
}
