package build

import (
	"path"
	"strings"

	"github.com/gobuffalo/plugins"
)

func FindBuilderFromArgs(args []string, plugs []plugins.Plugin) Builder {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		return FindBuilder(a, plugs)
	}
	return nil
}

func FindBuilder(name string, plugs []plugins.Plugin) Builder {
	// Find wraps the other cmd finders into a mega finder for cmds
	for _, p := range plugs {
		c, ok := p.(Builder)
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
