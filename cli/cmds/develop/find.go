package develop

import (
	"path"
	"strings"

	"github.com/gobuffalo/plugins"
)

func FindDeveloperFromArgs(args []string, plugs []plugins.Plugin) Developer {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		return FindDeveloper(a, plugs)
	}
	return nil
}

func FindDeveloper(name string, plugs []plugins.Plugin) Developer {
	// Find wraps the other cmd finders into a mega finder for cmds
	for _, p := range plugs {
		c, ok := p.(Developer)
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
