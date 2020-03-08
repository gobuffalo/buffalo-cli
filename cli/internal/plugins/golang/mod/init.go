package mod

import (
	"context"
	"os/exec"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
)

var _ plugins.Plugin = &Initer{}
var _ plugins.Needer = &Initer{}
var _ plugins.Scoper = &Initer{}

type Initer struct {
	pluginsFn plugins.Feeder
}

func (i *Initer) WithPlugins(f plugins.Feeder) {
	i.pluginsFn = f
}

func (i *Initer) ScopedPlugins() []plugins.Plugin {
	if i.pluginsFn == nil {
		return nil
	}

	var plugs []plugins.Plugin
	for _, p := range i.pluginsFn() {
		switch p.(type) {
		case Stderrer:
			plugs = append(plugs, p)
		case Stdiner:
			plugs = append(plugs, p)
		case Stdouter:
			plugs = append(plugs, p)
		case Requirer:
			plugs = append(plugs, p)
		case Replacer:
			plugs = append(plugs, p)
		}
	}

	return plugs
}

func (i Initer) PluginName() string {
	return "go/mod/init"
}

func (i *Initer) ModInit(ctx context.Context, root string, name string) error {
	plugs := i.ScopedPlugins()

	c := exec.CommandContext(ctx, "go", "mod", "init", name)
	c.Stdout = plugio.Stdout(plugs...)
	c.Stderr = plugio.Stderr(plugs...)
	c.Stdin = plugio.Stdin(plugs...)

	if err := c.Run(); err != nil {
		return err
	}
	return nil
}
