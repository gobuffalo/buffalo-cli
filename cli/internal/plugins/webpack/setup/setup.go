package setup

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/setup"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack/internal/scripts"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"
)

var _ plugcmd.Namer = &Setup{}
var _ plugins.Needer = &Setup{}
var _ plugins.Plugin = &Setup{}
var _ plugins.Scoper = &Setup{}
var _ setup.BeforeSetuper = &Setup{}

type Setup struct {
	pluginsFn plugins.Feeder
}

func (Setup) PluginName() string {
	return "webpack/setup"
}

func (Setup) CmdName() string {
	return "webpack"
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
		case Stdiner:
			plugs = append(plugs, p)
		case Stdouter:
			plugs = append(plugs, p)
		case Stderrer:
			plugs = append(plugs, p)
		}
	}

	return plugs
}

func (s *Setup) BeforeSetup(ctx context.Context, root string, args []string) error {
	tool, err := scripts.Tool(s, ctx, root)
	if err != nil {
		return err
	}

	// make sure that the node_modules folder is properly "installed"
	if _, err := os.Stat(filepath.Join(root, "node_modules")); err != nil {
		cmd := s.cmd(ctx, tool, "install")
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	cmd := s.cmd(ctx, scripts.WebpackBin(root))

	if _, err := scripts.ScriptFor(s, ctx, root, "build"); err == nil {
		cmd = s.cmd(ctx, tool, "run", "build")
	}

	return cmd.Run()
}

func (d *Setup) cmd(ctx context.Context, tool string, args ...string) *exec.Cmd {
	plugs := d.ScopedPlugins()

	c := exec.CommandContext(ctx, tool, args...)
	c.Stdin = plugio.Stdin(plugs...)
	c.Stdout = plugio.Stdout(plugs...)
	c.Stderr = plugio.Stderr(plugs...)
	return c
}
