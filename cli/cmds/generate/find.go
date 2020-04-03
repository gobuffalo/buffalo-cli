package generate

import (
	"path"
	"strings"

	"github.com/gobuffalo/plugins"
)

func FindGeneratorFromArgs(args []string, plugs []plugins.Plugin) Generator {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		return FindGenerator(a, plugs)
	}
	return nil
}

func FindGenerator(name string, plugs []plugins.Plugin) Generator {
	// Find wraps the other cmd finders into a mega finder for cmds
	for _, p := range plugs {
		c, ok := p.(Generator)
		if !ok {
			continue
		}
		if n, ok := c.(Namer); ok {
			if n.CmdName() == name {
				return c
			}
		}

		if n, ok := c.(Aliaser); ok {
			for _, a := range n.CmdAliases() {
				if a == name {
					return c
				}
			}
		}

		if name == path.Base(c.PluginName()) {
			return c
		}
	}
	return nil
}
