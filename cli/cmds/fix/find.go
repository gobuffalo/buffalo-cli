package fix

import (
	"path"
	"strings"

	"github.com/gobuffalo/plugins"
)

func FindFixerFromArgs(args []string, plugs []plugins.Plugin) Fixer {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		return FindFixer(a, plugs)
	}
	return nil
}

func FindFixer(name string, plugs []plugins.Plugin) Fixer {
	// Find wraps the other cmd finders into a mega finder for cmds
	for _, p := range plugs {
		c, ok := p.(Fixer)
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
